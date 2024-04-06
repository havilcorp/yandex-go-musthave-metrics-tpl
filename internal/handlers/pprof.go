package handlers

import (
	"net/http/pprof"

	"github.com/go-chi/chi"
)

type PPROFHandler struct{}

func NewPPROFHandler() *PPROFHandler {
	return &PPROFHandler{}
}

func (h *PPROFHandler) Register(router *chi.Mux) {
	// router.HandleFunc("/pprof/*", pprof.Index)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
}
