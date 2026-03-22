package handlers

import (
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Welcome to School API"))
	if err != nil {
		log.Fatal("Error writing to server:", err)
	}
}
