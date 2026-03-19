package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Major string `json:"major"`
}

func main() {
	port := ":3000"

	cert := "cert.pem"
	key := "key.pem"

	router := router.Router()

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
	secureMux := mw.SecurityHeaders(router)

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
