package main

import (
	"fmt"
	"log"
	"net/http"
	"wss-chat/handlers"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("starting...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error .env")
	}

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("GET /ws", handlers.Websocket)


	log.Fatal(http.ListenAndServe(":8081", nil))
}
