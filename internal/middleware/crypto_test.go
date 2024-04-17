// Package middleware мидлвар для шиврования и расшифрования запросов
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHashSHA256Middleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Error(err)
		}
	})
	r := httptest.NewRequest(http.MethodPost, "/ping", strings.NewReader("OK"))
	h := hmac.New(sha256.New, []byte("key"))
	_, err := h.Write([]byte("OK"))
	if err != nil {
		t.Error(err)
	}
	r.Header.Add("HashSHA256", hex.EncodeToString(h.Sum(nil)))
	rw := httptest.NewRecorder()
	HashSHA256Middleware("key")(testHandler).ServeHTTP(rw, r)
	res := rw.Result()
	assert.Equal(t, 200, res.StatusCode)
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()
}