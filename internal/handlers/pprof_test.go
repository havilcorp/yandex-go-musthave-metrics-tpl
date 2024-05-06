package handlers

import (
	"testing"

	"github.com/go-chi/chi"
)

func TestPPROFHandler_Register(t *testing.T) {
	r := chi.NewRouter()
	h := NewPPROFHandler()
	h.Register(r)
}
