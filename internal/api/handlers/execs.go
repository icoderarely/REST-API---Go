package handlers

import (
	"log"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("GET: Execs Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodPost:
		_, err := w.Write([]byte("POST: Execs Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodPut:
		_, err := w.Write([]byte("PUT: Execs Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodDelete:
		_, err := w.Write([]byte("DELETE: Execs Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	}
}
