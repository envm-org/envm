package middleware

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")

		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		w.Header().Set("X-Content-Type-Options", "nosniff")

		w.Header().Set("X-Frame-Options", "DENY")

		w.Header().Set("X-XSS-Protection", "1; mode=block")

		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
