package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func HashSHA256Middleware(key string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("HashSHA256") != "" {
				hexSha256, err := hex.DecodeString(r.Header.Get("HashSHA256"))
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				err = r.Body.Close()
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				h := hmac.New(sha256.New, []byte(key))
				h.Write(body)
				if !hmac.Equal(h.Sum(nil), hexSha256) {
					w.WriteHeader(http.StatusBadRequest)
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
