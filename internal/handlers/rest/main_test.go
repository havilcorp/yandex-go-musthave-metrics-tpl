package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func TestMainHandler_MainPageHandler(t *testing.T) {
	mainHandler := mocks.NewIMain(t)

	mainHandler.On("GetAllCounters", mock.Anything).Return(map[string]int64{"count": 1}, nil)
	mainHandler.On("GetAllGauge", mock.Anything).Return(map[string]float64{"gauge": 1.1}, nil)

	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MainPageHandler",
			args: args{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rw := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			h := NewMainHandler(mainHandler)
			h.MainPageHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Error(err)
				}
			}()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) == "" {
				t.Errorf("expected ABC got %v", string(data))
			}
		})
	}
}

func TestMainHandler_Register(t *testing.T) {
	r := chi.NewRouter()
	mainHandler := mocks.NewIMain(t)
	h := NewMainHandler(mainHandler)
	h.Register(r)
}
