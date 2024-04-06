package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/mocks"
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
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}
