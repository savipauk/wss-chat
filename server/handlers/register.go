package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

func Register(w http.ResponseWriter, r *http.Request) {
    var data RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


    log.Println(data)

	// w.Write([]byte("hi\n"))
}
