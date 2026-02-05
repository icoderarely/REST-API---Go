package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Major string `json:"major"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello Root Route")
	_, err := w.Write([]byte("Hello Root Route"))
	if err != nil {
		log.Fatal("Error writing to server:", err)
	}
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("GET: Teachers Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodPost:
		_, err := w.Write([]byte("POST: Teachers Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodPut:
		_, err := w.Write([]byte("PUT: Teachers Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodDelete:
		_, err := w.Write([]byte("DELETE: Teachers Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	}
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
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
	_, err := w.Write([]byte("Hello Execs Route"))
	if err != nil {
		log.Fatal("Error writing to server:", err)
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	port := ":3000"

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port[1:])

	if err := server.ListenAndServeTLS(cert, key); err != nil {
		log.Fatal("Error handling the server", err)
	}
}
