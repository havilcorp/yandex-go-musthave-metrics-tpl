// Package middleware мидлвар для проверки ip адреса
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTrustedSubnetMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	})
	t.Run("Good", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/ping", nil)
		r.Header.Add("X-Real-IP", "192.168.0.112")
		rw := httptest.NewRecorder()
		TrustedSubnetMiddleware("192.168.0.0/24")(testHandler).ServeHTTP(rw, r)
		res := rw.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer func() {
			if err := res.Body.Close(); err != nil {
				logrus.Info(err)
			}
		}()
	})
	t.Run("StatusForbidden", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/ping", nil)
		r.Header.Add("X-Real-IP", "192.168.1.112")
		rw := httptest.NewRecorder()
		TrustedSubnetMiddleware("192.168.0.0/24")(testHandler).ServeHTTP(rw, r)
		res := rw.Result()
		assert.Equal(t, http.StatusForbidden, res.StatusCode)
		defer func() {
			if err := res.Body.Close(); err != nil {
				logrus.Info(err)
			}
		}()
	})
	t.Run("StatusForbidden", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/ping", nil)
		rw := httptest.NewRecorder()
		TrustedSubnetMiddleware("192.168.0.0/24")(testHandler).ServeHTTP(rw, r)
		res := rw.Result()
		assert.Equal(t, http.StatusForbidden, res.StatusCode)
		defer func() {
			if err := res.Body.Close(); err != nil {
				logrus.Info(err)
			}
		}()
	})
}
