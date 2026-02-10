package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received Request from ResponseTime")
		start := time.Now()

		// custom responsewriter to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(w, r)

		// calculate the duration
		duration := time.Since(start)

		// log the request details
		fmt.Printf("Method: %s, URL: %s, Status %d, Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.Seconds())
		fmt.Println("Sent Response from ResponseTime middleware")
	})
}

// response writer
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
