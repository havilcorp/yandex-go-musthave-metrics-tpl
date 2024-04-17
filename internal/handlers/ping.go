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

// NewPingHandler инициализация
func NewPingHandler(pingService Pinger) *PingHandler {
	return &PingHandler{pingService: pingService}
}

// Register регистрация роутов
func (h *PingHandler) Register(router *chi.Mux) {
	router.Get("/ping", h.Ping)
}

// Ping хендлер для проверки подлкючения к базе данных
func (h *PingHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	if err := h.pingService.Ping(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}
