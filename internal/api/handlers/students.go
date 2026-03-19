package handlers

import (
	"log"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("GET: Students Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodPost:
		_, err := w.Write([]byte("POST: Students Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodPut:
		_, err := w.Write([]byte("PUT: Students Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodDelete:
		_, err := w.Write([]byte("DELETE: Students Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	}
	_, err := w.Write([]byte("Hello Students Route"))
	if err != nil {
		log.Fatal("Error writing to server:", err)
	}
}
