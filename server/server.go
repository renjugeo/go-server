package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/renjugeo/go-server/api"
	"github.com/renjugeo/go-server/config"
	"go.uber.org/zap"
)

type handler struct {
	// can embed router to server
	router *mux.Router
	cfg    *config.Configuration
	logger *zap.Logger
}

func NewServer(api *api.API, cfg *config.Configuration, logger *zap.Logger) *http.Server {
	h := &handler{
		router: mux.NewRouter(),
		cfg:    cfg,
		logger: logger,
	}
	h.setupRoutes()
	if api != nil {
		api.RegisterPaths(h.router)
	}
	return &http.Server{
		Handler:      h,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (s *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
