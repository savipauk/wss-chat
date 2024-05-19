package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wss-chat/handlers"
	"wss-chat/jwt"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("starting...")

    err := godotenv.Load()
    if err != nil {
        log.Fatal("error .env")
    }

    test := os.Getenv("test")
    fmt.Println(test)


    testToken, _ := jwt.Construct("johndoe")
    fmt.Println(testToken)
    fmt.Println()
    jwt.DecodeBase64(testToken)


	// http.HandleFunc("/", handleWebSocket)

	http.HandleFunc("GET /{$}", handlers.Root)
	http.HandleFunc("POST /register", handlers.Register)

    http.HandleFunc("GET /ws", connectWebSocket)

	// fileServer := http.FileServer(http.Dir("./handlers/"))
	// http.Handle("/handlers", http.StripPrefix("/handlers/", fileServer))

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func connectWebSocket (w http.ResponseWriter, r *http.Request) {
    handleWebSocket(w, r)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close()

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Recieved msg: %s\n", message)

		err = connection.WriteMessage(messageType, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
