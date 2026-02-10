package middleware

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content Security Policy (CSP)
		// Prevent XSS by only allowing scripts from self. Adjust as needed for external scripts/styles.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")

		// Strict Transport Security (HSTS)
		// Enforce HTTPS for 1 year (31536000 seconds) including subdomains
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Content-Type-Options
		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options
		// Prevent Clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// X-XSS-Protection
		// Enable XSS filtering in browser (though largely deprecated in modern browsers, still good for catch-all)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy
		// Control how much referrer information is sent
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
