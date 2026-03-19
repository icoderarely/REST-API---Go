package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	mw "restapi/internal/api/middlewares"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Major string `json:"major"`
}

type Teacher struct {
	ID        int
	FirstName string
	LastName  string
	Class     string
	Subject   string
}

var (
	teachers = make(map[int]Teacher)
	// mutex    = &sync.Mutex{}
	nextID = 1
)

// initialize some data
func init() {
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Algebra",
	}
	nextID++
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "113",
		Subject:   "Computer",
	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	fname := r.URL.Query().Get("first_name")
	lname := r.URL.Query().Get("last_name")

	teacherList := make([]Teacher, 0, len(teachers))

	if idStr == "" {
		for _, teacher := range teachers {
			// filter based on first_name/last_name query param
			if (fname == "" || teacher.FirstName == fname) && (lname == "" || teacher.LastName == lname) {
				teacherList = append(teacherList, teacher)
			}
		}
		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		id, _ := strconv.Atoi(idStr)
		teacher, exists := teachers[id]

		if !exists {
			http.Error(w, "ID not found", http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teacher)
	}
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
		// call get method handler function
		getTeachersHandler(w, r)
	// _, err := w.Write([]byte("GET: Teachers Route"))
	// if err != nil {
	// 	log.Fatal("Error writing to server:", err)
	// }
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
	// _, err := w.Write([]byte("Hello Execs Route"))
	// if err != nil {
	// 	log.Fatal("Error writing to server:", err)
	// }
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

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "applicaton/x-www-form-urlencode",
	// 	Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	// proper logical and efficient order of middlewares
	// secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	// secureMux := applyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)

	// for faster dev, will uncomment the rest middlewares later
	secureMux := mw.SecurityHeaders(mux)

	// create custom server
	server := &http.Server{
		Addr: port,
		// Handler:   middlewares.SecurityHeaders(mux),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port[1:])

	if err := server.ListenAndServeTLS(cert, key); err != nil {
		log.Fatal("Error handling the server", err)
	}
}

// Middleware is a function that wraps http.Handler with additional functionality
// type Middleware func(http.Handler) http.Handler

// func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
// 	for _, middleware := range middlewares {
// 		handler = middleware(handler)
// 	}
// 	return handler
// }
