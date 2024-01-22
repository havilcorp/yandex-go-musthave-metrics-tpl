package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type Want struct {
	method string
	code   int
	req    *http.Request
}

var storeTest = memory.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

func AddChiURLParams(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestUpdateCounterHandler(t *testing.T) {
	tests := []struct {
		name string
		want Want
	}{
		{
			name: "Bad value",
			want: Want{
				method: http.MethodPost,
				code:   400,
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/counter/Alloc/10f", nil), map[string]string{
					"name":  "Alloc",
					"value": "10f",
				}),
			},
		},
		{
			name: "Good",
			want: Want{
				method: http.MethodPost,
				code:   200,
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/counter/Alloc/10", nil), map[string]string{
					"name":  "Alloc",
					"value": "10",
				}),
			},
		},
	}
	SetStore(&storeTest)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			UpdateCounterHandler(w, test.want.req)
			res := w.Result()
			require.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestUpdateGaugeHandler(t *testing.T) {
	tests := []struct {
		name string
		want Want
	}{
		{
			name: "Bad value",
			want: Want{
				method: http.MethodPost,
				code:   400,
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/gauge/Alloc/10f", nil), map[string]string{
					"name":  "Alloc",
					"value": "10f",
				}),
			},
		},
		{
			name: "Good",
			want: Want{
				method: http.MethodPost,
				code:   200,
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/gauge/Alloc/10", nil), map[string]string{
					"name":  "Alloc",
					"value": "10",
				}),
			},
		},
	}
	SetStore(&storeTest)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			UpdateGaugeHandler(w, test.want.req)
			res := w.Result()
			require.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestGetCounterMetricHandler(t *testing.T) {
	tests := []struct {
		name string
		want Want
	}{
		{
			name: "Bad value",
			want: Want{
				method: http.MethodPost,
				code:   404,
				req: AddChiURLParams(httptest.NewRequest("GET", "/counter/Undefined", nil), map[string]string{
					"name": "Undefined",
				}),
			},
		},
		{
			name: "Good",
			want: Want{
				method: http.MethodPost,
				code:   200,
				req: AddChiURLParams(httptest.NewRequest("GET", "/counter/Alloc", nil), map[string]string{
					"name": "Alloc",
				}),
			},
		},
	}
	SetStore(&storeTest)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			GetCounterMetricHandler(w, test.want.req)
			res := w.Result()
			require.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestGetGaugeMetricHandler(t *testing.T) {
	tests := []struct {
		name string
		want Want
	}{
		{
			name: "Bad value",
			want: Want{
				method: http.MethodPost,
				code:   404,
				req: AddChiURLParams(httptest.NewRequest("GET", "/gauge/Undefined", nil), map[string]string{
					"name": "Undefined",
				}),
			},
		},
		{
			name: "Good",
			want: Want{
				method: http.MethodPost,
				code:   200,
				req: AddChiURLParams(httptest.NewRequest("GET", "/gauge/Alloc", nil), map[string]string{
					"name": "Alloc",
				}),
			},
		},
	}
	SetStore(&storeTest)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			GetGaugeMetricHandler(w, test.want.req)
			res := w.Result()
			require.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestMainPageHandler(t *testing.T) {
	t.Run("MainPage", func(t *testing.T) {
		w := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		MainPageHandler(w, request)
		res := w.Result()
		require.Equal(t, 200, res.StatusCode)
		defer res.Body.Close()
	})
}

func TestBadRequestHandler(t *testing.T) {
	t.Run("BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		BadRequestHandler(w, request)
		res := w.Result()
		require.Equal(t, 400, res.StatusCode)
		defer res.Body.Close()
	})
}
