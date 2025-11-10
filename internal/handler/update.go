package handler

import (
	"net/http"
	"strings"

	"github.com/RuslanBoo/metrics/internal/service"
)

type Handler struct {
	svc *service.MetricService
}

func New(svc *service.MetricService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/update/")
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}

	metricType, name, valueStr := parts[0], parts[1], parts[2]

	err := h.svc.UpdateMetric(metricType, name, valueStr)
	if err != nil {
		switch err {
		case service.ErrUnknownMetricType, service.ErrInvalidValue:
			http.Error(w, "", http.StatusBadRequest)
		default:
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
