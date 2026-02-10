package middlewares

import "net/http"

// basic middleware skeleton
// func securityHeaders(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// next.ServeHTTP(w, r)
// 	})
// }

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic security headers
		w.Header().Set("X-DNS-Prefetch-Control", "off")                                              // DNS prefetching attacks, prevents automatic DNS lookups
		w.Header().Set("X-Frame-Options", "DENY")                                                    // Clickjacking attacks; prevents embedding in iframes
		w.Header().Set("X-XSS-Protection", "1; mode=block")                                          // Cross-site scripting (XSS) attacks; enables browser XSS filtering
		w.Header().Set("X-Content-Type-Options", "nosniff")                                          // MIME sniffing attacks; prevents browser from guessing content types
		w.Header().Set("Strict-Transport-Security", "max-age=630772000; includeSubDomains; preload") // Forces HTTPS connections; prevents man-in-the-middle attacks
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")                         // Information disclosure; controls referrer header information
		w.Header().Set("X-Powered-By", "Django")                                                     // Fingerprinting attacks; obscures server technology

		// Content Security Policy - prevents XSS, code injection, clickjacking
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; media-src 'self'; object-src 'none'; child-src 'none'; worker-src 'self'; frame-ancestors 'none'; form-action 'self'; base-uri 'self'; manifest-src 'self'")

		// Permissions Policy - prevents unauthorized access to browser features
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), accelerometer=(), gyroscope=(), magnetometer=(), payment=(), usb=(), interest-cohort=()")

		// Cross-Origin policies
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp") // Prevents loading cross-origin resources without explicit permission
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")    // Isolates browsing context group; prevents cross-origin attacks
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")  // Prevents cross-origin resource loading; reduces side-channel attacks

		// Cache control for sensitive data
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate") // Prevents sensitive data caching; avoids data exposure
		w.Header().Set("Pragma", "no-cache")                                                     // Legacy cache control; prevents HTTP/1.0 caching
		w.Header().Set("Expires", "0")                                                           // Forces immediate expiration; prevents cached sensitive data

		// Additional security headers
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none") // Prevents Flash/PDF cross-domain policy attacks
		w.Header().Set("Expect-CT", "max-age=86400, enforce")       // Certificate Transparency monitoring; detects mis-issued certificates
		w.Header().Set("X-Download-Options", "noopen")              // Prevents IE from executing downloads in site's context

		next.ServeHTTP(w, r)
	})
}
