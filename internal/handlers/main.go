package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type IMain interface {
	GetAllCounters(ctx context.Context) (map[string]int64, error)
	GetAllGauge(ctx context.Context) (map[string]float64, error)
}

type MainHandler struct {
	main IMain
}

func NewMainHandler(main IMain) *MainHandler {
	return &MainHandler{
		main: main,
	}
}

func (h *MainHandler) Register(router *chi.Mux) {
	router.Get("/", h.MainPageHandler)
}

func (h *MainHandler) MainPageHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	liCounter := ""
	listCounters, err := h.main.GetAllCounters(r.Context())
	if err != nil {
		logrus.Errorf("GetAllCounters %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	for key, item := range listCounters {
		liCounter += fmt.Sprintf("<li>%s: %d</li>", key, item)
	}
	listGauges, err := h.main.GetAllGauge(r.Context())
	if err != nil {
		logrus.Errorf("GetAllGauge %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	liGauge := ""
	for key, item := range listGauges {
		liGauge += fmt.Sprintf("<li>%s: %f</li>", key, item)
	}
	html := fmt.Sprintf(`
	<html>
		<body>
			<br/>
			Counters
			<br/>
			<ul>%s</ul>
			<br/>
			Gauges
			<br/>
			<ul>%s</ul>
		</body>
	</html>
	`, liCounter, liGauge)
	_, err = rw.Write([]byte(html))
	if err != nil {
		logrus.Error(err)
	}
}
