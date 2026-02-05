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

			// Access request details
			fmt.Println("Body:", r.Body)                          // &{0x1400000e1e0 <nil> <nil> false true {{} {0 0}} true false false 0x1009dba40}
			fmt.Println("Form:", r.Form)                          // map[]
			fmt.Println("Header:", r.Header)                      //  map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Connection:[keep-alive] Content-Length:[63] Content-Type:[application/json] Postman-Token:[92a452d3-bf22-4baf-80d0-dabd65b9fb58] User-Agent:[PostmanRuntime/7.51.1]]
			fmt.Println("Context:", r.Context())                  // context.Background.WithValue(net/http context value http-server, *http.Server).WithValue(net/http context value local-addr, [::1]:3000).WithCancel.WithCancel
			fmt.Println("Content Len:", r.ContentLength)          // 63
			fmt.Println("Host:", r.Host)                          // localhost:3000
			fmt.Println("Method:", r.Method)                      // POST
			fmt.Println("Protocol:", r.Proto)                     // HTTP/1.1
			fmt.Println("Remote addr:", r.RemoteAddr)             // [::]:58871
			fmt.Println("Req URI:", r.RequestURI)                 // /students
			fmt.Println("TLS:", r.TLS)                            // <nil>
			fmt.Println("Trailer:", r.Trailer)                    // map[]
			fmt.Println("Transfer Encoding:", r.TransferEncoding) // []
			fmt.Println("URL:", r.URL)                            // /students
			fmt.Println("User agent:", r.UserAgent())             // PostmanRuntime/7.51.1
			fmt.Println("Port in URL:", r.URL.Port())             //
			fmt.Println("Scheme in URL:", r.URL.Scheme)           //

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
