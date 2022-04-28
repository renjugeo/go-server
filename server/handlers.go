package server

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *handler) setupRoutes() {
	s.router.HandleFunc("/", s.handleIndex()).Methods(http.MethodGet)
	s.router.HandleFunc(s.cfg.HealthUri, s.handleHealthCheck()).Methods(http.MethodGet)
	s.router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
}

func (s *handler) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body>Example Go Server</body></html>")
	}
}

func (s *handler) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}
}
