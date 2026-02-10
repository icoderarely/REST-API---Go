package middlewares

import (
	"fmt"
	"net/http"
	"slices"
)

// Allowed Origins
var allowedOrigins = []string{
	"https://my-origin-url.com",
	"https://localhost:3000",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println(origin)

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !isOriginAllowed(origin) {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")                         // Restricts allowed HTTP methods; prevents unauthorized operations
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin") // Controls allowed request headers; prevents header injection attacks
		w.Header().Set("Access-Control-Allow-Credentials", "true")                                                      // Controls credential sharing; manages authentication in cross-origin requests
		w.Header().Set("Access-Control-Max-Age", "3600")                                                                // Reduces preflight requests; improves performance while maintaining security

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	// for _, allowedOrigin := range allowedOrigins {
	// 	if origin == allowedOrigin {
	// 		return true
	// 	}
	// }
	// return false
	return slices.Contains(allowedOrigins, origin)
}
