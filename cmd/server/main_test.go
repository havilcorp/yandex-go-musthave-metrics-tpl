package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/handlers"
	"github.com/stretchr/testify/require"
)

type Want struct {
	method string
	code   int
	url    string
}

func TestUpdateCounterHandler(t *testing.T) {
	tests := []struct {
		name string
		want Want
	}{
		{
			name: "Empty type",
			want: Want{
				method: http.MethodPost,
				code:   404,
				url:    "/update",
			},
		},
		{
			name: "Bad type",
			want: Want{
				method: http.MethodPost,
				code:   404,
				url:    "/update/type/name/10",
			},
		},
		{
			name: "Bad value",
			want: Want{
				method: http.MethodPost,
				code:   400,
				url:    "/update/counter/name/10f",
			},
		},
		{
			name: "Bad method",
			want: Want{
				method: http.MethodGet,
				code:   404,
				url:    "/update/counter/name/10",
			},
		},
		{
			name: "Good",
			want: Want{
				method: http.MethodPost,
				code:   200,
				url:    "/update/counter/name/10",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.want.method, test.want.url, nil)
			w := httptest.NewRecorder()
			handlers.UpdateCounterHandler(w, request)
			res := w.Result()
			require.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}

}

func TestUpdateGuageHandler(t *testing.T) {
	tests := []struct {
		name string
		want Want
	}{
		{
			name: "Empty type",
			want: Want{
				method: http.MethodPost,
				code:   404,
				url:    "/update",
			},
		},
		{
			name: "Bad type",
			want: Want{
				method: http.MethodPost,
				code:   404,
				url:    "/update/type/name/10",
			},
		},
		{
			name: "Bad value",
			want: Want{
				method: http.MethodPost,
				code:   400,
				url:    "/update/gauge/name/10f",
			},
		},
		{
			name: "Bad method",
			want: Want{
				method: http.MethodGet,
				code:   404,
				url:    "/update/gauge/name/10",
			},
		},
		{
			name: "Good",
			want: Want{
				method: http.MethodPost,
				code:   200,
				url:    "/update/gauge/name/10",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.want.method, test.want.url, nil)
			w := httptest.NewRecorder()

			handlers.UpdateGaugeHandler(w, request)
			res := w.Result()
			require.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}

}
