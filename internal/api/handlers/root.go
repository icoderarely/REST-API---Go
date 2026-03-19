package handlers

import (
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello Root Route"))
	if err != nil {
		log.Fatal("Error writing to server:", err)
	}
}
