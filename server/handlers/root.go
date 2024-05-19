package handlers

import (
	"log"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Write([]byte("hi\n"))
}
