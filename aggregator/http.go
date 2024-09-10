package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"tolling/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPMetrictHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func newHTTPMetrictHandler(reqName string) *HTTPMetrictHandler {
	return &HTTPMetrictHandler{
		reqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_total", reqName),
			Name:      "requests",
		}),
		reqLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: fmt.Sprintf("http_%s_latency", reqName),
			Name:      "requests",
			Buckets:   []float64{0.1, 0.5, 1},
		}),
	}
}

func (h *HTTPMetrictHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()

			logrus.WithFields(logrus.Fields{
				"method":  r.Method,
				"path":    r.URL.Path,
				"latency": fmt.Sprintf("%.2fs", latency),
			}).Info("request processed")

			h.reqLatency.Observe(latency)
		}(time.Now())

		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not supported"})
			return
		}

		q, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}

		OBUID, err := strconv.Atoi(q[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}

		invoice, err := svc.CalculateInvoice(OBUID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not supported"})
			return
		}

		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
