package handlers

import (
	"testing"

	"github.com/go-chi/chi"
)

func TestPPROFHandler_Register(t *testing.T) {
	r := chi.NewRouter()
	t.Run("Register", func(t *testing.T) {
		h := NewPPROFHandler()
		h.Register(r)
	})
}
