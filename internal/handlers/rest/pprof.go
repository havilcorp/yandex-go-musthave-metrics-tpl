// Package rest роуты сервера
package rest

import (
	"net/http/pprof"

	"github.com/go-chi/chi"
)

type PPROFHandler struct{}

// NewPPROFHandler инициализация хендлера для сбора статистики
func NewPPROFHandler() *PPROFHandler {
	return &PPROFHandler{}
}

// Register регистрация роутов
func (h *PPROFHandler) Register(router *chi.Mux) {
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
}
