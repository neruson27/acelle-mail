package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (hh HealthHandler) NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/ping", handler(hh.health))
	return r
}

func (hh HealthHandler) health(w http.ResponseWriter, _ *http.Request) error {
	_ = json.NewEncoder(w).Encode("pong")
	return nil
}
