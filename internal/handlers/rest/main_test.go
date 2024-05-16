package rest

import (
	"errors"
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

	r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rw := httptest.NewRecorder()

	h := NewMainHandler(mainHandler)
	h.MainPageHandler(rw, r)
	res := rw.Result()
	assert.Equal(t, 200, res.StatusCode)
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) == "" {
		t.Errorf("expected got %v", string(data))
	}
}

func TestMainHandler_MainPageHandlerError(t *testing.T) {
	mainHandler := mocks.NewIMain(t)

	mainHandler.On("GetAllCounters", mock.Anything).Return(map[string]int64{"count": 1}, errors.New(""))
	mainHandler.On("GetAllGauge", mock.Anything).Return(map[string]float64{"gauge": 1.1}, errors.New(""))

	r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rw := httptest.NewRecorder()

	h := NewMainHandler(mainHandler)
	h.MainPageHandler(rw, r)
	res := rw.Result()
	assert.Equal(t, 500, res.StatusCode)
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Info(err)
		}
	}()
}

func TestMainHandler_Register(t *testing.T) {
	r := chi.NewRouter()
	mainHandler := mocks.NewIMain(t)
	h := NewMainHandler(mainHandler)
	h.Register(r)
}
