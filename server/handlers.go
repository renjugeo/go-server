package server

import (
	"fmt"
	"net/http"
)

func (s *handler) setupRoutes() {
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.router.HandleFunc("/health", s.handleHealthCheck()).Methods("GET")
}

func (s *handler) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body>Hello world</body></html>")
	}
}

func (s *handler) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}
}
