package middleware

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   string
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	rw.body += string(p)
	return rw.ResponseWriter.Write(p)
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming Request: method %s, endpoint: %s", r.Method, r.URL)
		log.Printf("Headers: %+v, Query Parameters: %+v", r.Header, r.URL.Query())

		rw := &responseWriter{w, http.StatusOK, ""}

		h.ServeHTTP(rw, r)

		log.Printf("Outgoing Response Headers: %+v", w.Header())
		log.Printf("status: %v", rw.status)
		if rw.status != http.StatusOK {
			log.Printf("error: %s", rw.body)
		}
	})
}

func BasicAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, pass, ok := r.BasicAuth()
		if !ok || !checkCredentials(username, pass) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// just mock, here could be checking in storage
func checkCredentials(username, password string) bool {
	return username == "user" && password == "password"
}
