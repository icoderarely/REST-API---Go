package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello Root Route")
		_, err := w.Write([]byte("Hello Root Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, err := w.Write([]byte("Hello GET Method on Teachers Route"))
			if err != nil {
				log.Fatal("Error writing to server:", err)
			}
			return
		}
		_, err := w.Write([]byte("Hello Teachers Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello Execs Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	})

	fmt.Println("Server is running on port:", port[1:])

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error handling the server", err)
	}
}
