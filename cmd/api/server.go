package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Major string `json:"major"`
}

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
			// parse form data (necessary for x-www-form-urlencoded) -> returns a map, with values in slice
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			fmt.Println("Form:", r.Form)
			// extract these values -> prepare response data
			response := make(map[string]interface{})
			for key, value := range r.Form {
				response[key] = value[0] // value[0] if we know we are getting only one value in form
			}
			fmt.Println("Response Map:", response)

			// RAW body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			defer r.Body.Close()

			fmt.Println("RAW Body:", string(body))

			// If you expect json data, then unmarshal it
			var userInstance User
			if err = json.Unmarshal(body, &userInstance); err != nil {
				return
			}
			fmt.Println(userInstance)

			_, err = w.Write([]byte("POST: Students Route"))
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
