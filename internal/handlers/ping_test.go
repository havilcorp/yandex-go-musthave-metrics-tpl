package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler_CheckDBHandler(t *testing.T) {
	pinger := mocks.NewPinger(t)
	pinger.On("Ping").Return(nil)

	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Ping",
			args: args{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rw := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			h := NewPingHandler(pinger)
			h.Ping(rw, r)
			res := rw.Result()
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Error(err)
				}
			}()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}

func TestPingHandler_Register(t *testing.T) {
	r := chi.NewRouter()
	pinger := mocks.NewPinger(t)
	h := NewPingHandler(pinger)
	h.Register(r)
}
