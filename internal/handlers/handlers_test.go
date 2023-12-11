package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

type Want struct {
	method string
	code   int
	req    *http.Request
}

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
