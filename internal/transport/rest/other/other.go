package other

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/transport/rest"
	"github.com/sirupsen/logrus"
)

type handler struct {
	store storage.IStorage
}

func NewHandler(store storage.IStorage) rest.Handler {
	return &handler{store: store}
}

func (h *handler) Register(router *chi.Mux) {
	router.Get("/", h.MainPageHandler)
	router.Get("/ping", h.CheckDBHandler)
}

func (h *handler) MainPageHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	liCounter := ""
	for key, item := range h.store.GetAllCounters(r.Context()) {
		liCounter += fmt.Sprintf("<li>%s: %d</li>", key, item)
	}
	liGauge := ""
	for key, item := range h.store.GetAllGauge(r.Context()) {
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
	_, err := rw.Write([]byte(html))
	if err != nil {
		logrus.Error(err)
	}
}

func (h *handler) CheckDBHandler(rw http.ResponseWriter, r *http.Request) {
	if err := h.store.Ping(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}
