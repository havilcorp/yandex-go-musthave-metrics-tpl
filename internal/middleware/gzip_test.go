// Package middleware мидлвар для сжатия и разжатия запросов
package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	})
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write([]byte("OK"))
	if err != nil {
		t.Error(err)
		return
	}
	if err = zb.Close(); err != nil {
		t.Error(err)
		return
	}
	r := httptest.NewRequest(http.MethodPost, "/ping", buf)
	r.Header.Add("Accept-Encoding", "gzip")
	r.Header.Add("Content-Encoding", "gzip")
	rw := httptest.NewRecorder()
	lm := GzipMiddleware(testHandler)
	lm.ServeHTTP(rw, r)
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()
}
