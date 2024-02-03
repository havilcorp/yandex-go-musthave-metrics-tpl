package rest

import (
	"github.com/go-chi/chi"
)

type Handler interface {
	Register(router *chi.Mux)
}

// var store storage.IStorage

// func SetStore(s storage.IStorage) {
// 	store = s
// }

// func CheckDBHandler(rw http.ResponseWriter, r *http.Request) {
// 	if err := store.Ping(); err != nil {
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	} else {
// 		rw.WriteHeader(http.StatusOK)
// 	}
// }

// func UpdateHandler(rw http.ResponseWriter, r *http.Request) {
// 	var req models.MetricsRequest
// 	rw.Header().Set("Content-Type", "application/json")
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&req); err != nil {
// 		logrus.Info(err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	if req.MType == models.TypeMetricsCounter {
// 		if err := store.AddCounter(r.Context(), req.ID, *req.Delta); err != nil {
// 			logrus.Info(err)
// 			rw.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		if val, ok := store.GetCounter(r.Context(), req.ID); ok {
// 			rw.WriteHeader(http.StatusOK)
// 			resp := models.MetricsRequest{
// 				ID:    req.ID,
// 				MType: req.MType,
// 				Delta: &val,
// 			}
// 			enc := json.NewEncoder(rw)
// 			if err := enc.Encode(resp); err != nil {
// 				logrus.Info(err)
// 				return
// 			}
// 		} else {
// 			rw.WriteHeader(http.StatusNotFound)
// 		}
// 	}
// 	if req.MType == models.TypeMetricsGauge {
// 		if err := store.AddGauge(r.Context(), req.ID, *req.Value); err != nil {
// 			logrus.Info(err)
// 			rw.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		if val, ok := store.GetGauge(r.Context(), req.ID); ok {
// 			rw.WriteHeader(http.StatusOK)
// 			resp := models.MetricsRequest{
// 				ID:    req.ID,
// 				MType: req.MType,
// 				Value: &val,
// 			}
// 			enc := json.NewEncoder(rw)
// 			if err := enc.Encode(resp); err != nil {
// 				logrus.Info(err)
// 				return
// 			}
// 		} else {
// 			rw.WriteHeader(http.StatusNotFound)
// 		}
// 	}

// }

// func UpdateBulkHandler(rw http.ResponseWriter, r *http.Request) {
// 	metrics := make([]models.MetricsRequest, 0)
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&metrics); err != nil {
// 		logrus.Info(err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	gauge := make([]models.GaugeModel, 0)
// 	counter := make([]models.CounterModel, 0)
// 	for _, m := range metrics {
// 		if m.MType == models.TypeMetricsGauge {
// 			gauge = append(gauge, models.GaugeModel{Key: m.ID, Value: *m.Value})
// 		} else if m.MType == models.TypeMetricsCounter {
// 			counter = append(counter, models.CounterModel{Key: m.ID, Value: *m.Delta})
// 		}
// 	}
// 	if err := store.AddGaugeBulk(r.Context(), gauge); err != nil {
// 		logrus.Info(err)
// 	}
// 	if err := store.AddCounterBulk(r.Context(), counter); err != nil {
// 		logrus.Info(err)
// 	}
// }

// func UpdateCounterHandler(rw http.ResponseWriter, r *http.Request) {

// 	marketName := chi.URLParam(r, "name")
// 	marketVal := chi.URLParam(r, "value")

// 	marketValInt64, err := strconv.ParseInt(marketVal, 0, 64)
// 	if err != nil {
// 		logrus.Info(err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	if err := store.AddCounter(r.Context(), marketName, marketValInt64); err != nil {
// 		logrus.Info(err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	rw.WriteHeader(http.StatusOK)
// }

// func UpdateGaugeHandler(rw http.ResponseWriter, r *http.Request) {

// 	marketName := chi.URLParam(r, "name")
// 	marketVal := chi.URLParam(r, "value")

// 	marketValFloat64, err := strconv.ParseFloat(marketVal, 64)
// 	if err != nil {
// 		logrus.Info(err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	if err := store.AddGauge(r.Context(), marketName, marketValFloat64); err != nil {
// 		logrus.Info(err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	rw.WriteHeader(http.StatusOK)

// }

// func GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
// 	var req models.MetricsRequest
// 	rw.Header().Set("Content-Type", "application/json")
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&req); err != nil {
// 		logrus.Info(err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	if req.MType == models.TypeMetricsCounter {
// 		if val, ok := store.GetCounter(r.Context(), req.ID); ok {
// 			rw.WriteHeader(http.StatusOK)
// 			resp := models.MetricsRequest{
// 				ID:    req.ID,
// 				MType: req.MType,
// 				Delta: &val,
// 			}
// 			enc := json.NewEncoder(rw)
// 			if err := enc.Encode(resp); err != nil {
// 				logrus.Info(err)
// 				return
// 			}
// 		} else {
// 			rw.WriteHeader(http.StatusNotFound)
// 		}
// 	}
// 	if req.MType == models.TypeMetricsGauge {
// 		if val, ok := store.GetGauge(r.Context(), req.ID); ok {
// 			rw.WriteHeader(http.StatusOK)
// 			resp := models.MetricsRequest{
// 				ID:    req.ID,
// 				MType: req.MType,
// 				Value: &val,
// 			}
// 			enc := json.NewEncoder(rw)
// 			if err := enc.Encode(resp); err != nil {
// 				logrus.Info(err)
// 				return
// 			}
// 		} else {
// 			rw.WriteHeader(http.StatusNotFound)
// 		}
// 	}
// }

// func GetCounterMetricHandler(rw http.ResponseWriter, r *http.Request) {
// 	marketName := chi.URLParam(r, "name")
// 	if val, ok := store.GetCounter(r.Context(), marketName); ok {
// 		rw.WriteHeader(http.StatusOK)
// 		_, err := rw.Write([]byte(fmt.Sprintf("%d", val)))
// 		if err != nil {
// 			logrus.Info(err)
// 		}
// 	} else {
// 		rw.WriteHeader(http.StatusNotFound)
// 	}
// }

// func GetGaugeMetricHandler(rw http.ResponseWriter, r *http.Request) {
// 	marketName := chi.URLParam(r, "name")
// 	if val, ok := store.GetGauge(r.Context(), marketName); ok {
// 		rw.WriteHeader(http.StatusOK)
// 		_, err := rw.Write([]byte(fmt.Sprintf("%g", val)))
// 		if err != nil {
// 			logrus.Info(err)
// 		}
// 	} else {
// 		rw.WriteHeader(http.StatusNotFound)
// 	}
// }

// func MainPageHandler(rw http.ResponseWriter, r *http.Request) {
// 	rw.Header().Set("Content-Type", "text/html")
// 	liCounter := ""
// 	for key, item := range store.GetAllCounters(r.Context()) {
// 		liCounter += fmt.Sprintf("<li>%s: %d</li>", key, item)
// 	}
// 	liGauge := ""
// 	for key, item := range store.GetAllGauge(r.Context()) {
// 		liGauge += fmt.Sprintf("<li>%s: %f</li>", key, item)
// 	}
// 	html := fmt.Sprintf(`
// 	<html>
// 		<body>
// 			<br/>
// 			Counters
// 			<br/>
// 			<ul>%s</ul>
// 			<br/>
// 			Gauges
// 			<br/>
// 			<ul>%s</ul>
// 		</body>
// 	</html>
// 	`, liCounter, liGauge)
// 	_, err := rw.Write([]byte(html))
// 	if err != nil {
// 		logrus.Info(err)
// 	}
// }

// func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
// 	rw.WriteHeader(http.StatusBadRequest)
// }

// package handlers

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi"
// 	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memory"
// 	"github.com/stretchr/testify/require"
// )

// type Want struct {
// 	method string
// 	code   int
// 	req    *http.Request
// }

// var storeTest = memory.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

// func AddChiURLParams(r *http.Request, params map[string]string) *http.Request {
// 	ctx := chi.NewRouteContext()
// 	for k, v := range params {
// 		ctx.URLParams.Add(k, v)
// 	}
// 	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
// }

// func TestUpdateCounterHandler(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want Want
// 	}{
// 		{
// 			name: "Bad value",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   400,
// 				req: AddChiURLParams(httptest.NewRequest("POST", "/update/counter/Alloc/10f", nil), map[string]string{
// 					"name":  "Alloc",
// 					"value": "10f",
// 				}),
// 			},
// 		},
// 		{
// 			name: "Good",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   200,
// 				req: AddChiURLParams(httptest.NewRequest("POST", "/update/counter/Alloc/10", nil), map[string]string{
// 					"name":  "Alloc",
// 					"value": "10",
// 				}),
// 			},
// 		},
// 	}
// 	SetStore(&storeTest)
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			UpdateCounterHandler(w, test.want.req)
// 			res := w.Result()
// 			require.Equal(t, test.want.code, res.StatusCode)
// 			defer res.Body.Close()
// 		})
// 	}
// }

// func TestUpdateGaugeHandler(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want Want
// 	}{
// 		{
// 			name: "Bad value",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   400,
// 				req: AddChiURLParams(httptest.NewRequest("POST", "/update/gauge/Alloc/10f", nil), map[string]string{
// 					"name":  "Alloc",
// 					"value": "10f",
// 				}),
// 			},
// 		},
// 		{
// 			name: "Good",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   200,
// 				req: AddChiURLParams(httptest.NewRequest("POST", "/update/gauge/Alloc/10", nil), map[string]string{
// 					"name":  "Alloc",
// 					"value": "10",
// 				}),
// 			},
// 		},
// 	}
// 	SetStore(&storeTest)
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			UpdateGaugeHandler(w, test.want.req)
// 			res := w.Result()
// 			require.Equal(t, test.want.code, res.StatusCode)
// 			defer res.Body.Close()
// 		})
// 	}
// }

// func TestGetCounterMetricHandler(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want Want
// 	}{
// 		{
// 			name: "Bad value",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   404,
// 				req: AddChiURLParams(httptest.NewRequest("GET", "/counter/Undefined", nil), map[string]string{
// 					"name": "Undefined",
// 				}),
// 			},
// 		},
// 		{
// 			name: "Good",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   200,
// 				req: AddChiURLParams(httptest.NewRequest("GET", "/counter/Alloc", nil), map[string]string{
// 					"name": "Alloc",
// 				}),
// 			},
// 		},
// 	}
// 	SetStore(&storeTest)
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			GetCounterMetricHandler(w, test.want.req)
// 			res := w.Result()
// 			require.Equal(t, test.want.code, res.StatusCode)
// 			defer res.Body.Close()
// 		})
// 	}
// }

// func TestGetGaugeMetricHandler(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want Want
// 	}{
// 		{
// 			name: "Bad value",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   404,
// 				req: AddChiURLParams(httptest.NewRequest("GET", "/gauge/Undefined", nil), map[string]string{
// 					"name": "Undefined",
// 				}),
// 			},
// 		},
// 		{
// 			name: "Good",
// 			want: Want{
// 				method: http.MethodPost,
// 				code:   200,
// 				req: AddChiURLParams(httptest.NewRequest("GET", "/gauge/Alloc", nil), map[string]string{
// 					"name": "Alloc",
// 				}),
// 			},
// 		},
// 	}
// 	SetStore(&storeTest)
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			GetGaugeMetricHandler(w, test.want.req)
// 			res := w.Result()
// 			require.Equal(t, test.want.code, res.StatusCode)
// 			defer res.Body.Close()
// 		})
// 	}
// }

// func TestMainPageHandler(t *testing.T) {
// 	t.Run("MainPage", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		request := httptest.NewRequest(http.MethodGet, "/", nil)
// 		MainPageHandler(w, request)
// 		res := w.Result()
// 		require.Equal(t, 200, res.StatusCode)
// 		defer res.Body.Close()
// 	})
// }

// func TestBadRequestHandler(t *testing.T) {
// 	t.Run("BadRequest", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		request := httptest.NewRequest(http.MethodGet, "/", nil)
// 		BadRequestHandler(w, request)
// 		res := w.Result()
// 		require.Equal(t, 400, res.StatusCode)
// 		defer res.Body.Close()
// 	})
// }
