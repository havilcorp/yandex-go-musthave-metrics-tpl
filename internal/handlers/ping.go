package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

//go:generate mockery --name Pinger
type Pinger interface {
	Ping() error
}

type PingHandler struct {
	pingService Pinger
}

func NewPingHandler(pingService Pinger) *PingHandler {
	return &PingHandler{pingService: pingService}
}

func (h *PingHandler) Register(router *chi.Mux) {
	router.Get("/ping", h.Ping)
}

func (h *PingHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	if err := h.pingService.Ping(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}
