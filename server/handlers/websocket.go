package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"wss-chat/jwt"
	"wss-chat/session"

	"github.com/gorilla/websocket"
)

type WSConnection struct {
	Connection *websocket.Conn
	Token      string
}

type Message struct {
	Message string `json:"message"`
	Token   string `token:"token"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections []WSConnection

func Websocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	// log.Print("token: ", token)

	// Read session file and check if token exists
	verified, _ := jwt.Verify(token)
	if !verified {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !session.CheckIfExists(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("erorr: %v", err)
		return
	}
	defer connection.Close()

	connections = append(connections, WSConnection{Connection: connection, Token: token})

	var wg sync.WaitGroup
	wg.Add(1)
	go DistributeMessages(connection, token)
	wg.Wait()
}

func DistributeMessages(connection *websocket.Conn, token string) {
	for {
		messageType, messageBytes, err := connection.ReadMessage()

		if !session.CheckIfExists(token) {
			log.Printf("Couldn't find session for token %s", token)
			connection.Close()
			return
		}

		if err != nil {
			log.Printf("Connection end: %v,\tToken: %s", err, token)
			return
		}

		// Message is type
		// { message: message, token: token }
		// SEND MESSAGE TO ALL CLIENTS WHICH ARE NOT ME (CHECK BY TOKEN (dok stvoris konekciju napravi si array *websocket.Conn))
		var message Message
		err = json.Unmarshal(messageBytes, &message)

		if err != nil {
			log.Printf("Error parsing message: %v", err)
			return
		}

		fmt.Printf("Recieved message (token): %s (%s)\n", message.Message, message.Token)

		for _, conn := range connections {
			if conn.Token == message.Token && conn.Connection == connection {
				continue
			}

			log.Printf("Sending message %s to %s\n", message.Message, conn.Token)
			err = conn.Connection.WriteMessage(messageType, []byte(message.Message))
			if err != nil {
				log.Printf("error while sendind msg (maybe client disconnected?): %v", err)
				continue
			}
		}
	}
}
