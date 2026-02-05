package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
	// teachers/{id} - path params
	fmt.Println(r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	userID := strings.TrimSuffix(path, "/")
	fmt.Println("The ID is:", userID)

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

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server is running on port:", port[1:])

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error handling the server", err)
	}
}
