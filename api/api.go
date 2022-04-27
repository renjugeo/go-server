package api

import (
	"github.com/gorilla/mux"
	v1 "github.com/renjugeo/go-server/api/v1"
	"github.com/renjugeo/go-server/config"
	"go.uber.org/zap"
)

type API struct {
	v1 *v1.API
}

func NewAPI(cfg *config.Configuration, logger *zap.Logger) *API {
	return &API{
		v1: v1.NewV1API(cfg, logger),
	}
}

func (api *API) RegisterPaths(r *mux.Router) {
	api.v1.RegisterPaths(r)
}
