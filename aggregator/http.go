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

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code int
	Err  error
}

// Error implements the error interface
func (e APIError) Error() string {
	return e.Err.Error()
}

type HTTPMetrictHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Error()})
				return
			}

			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func newHTTPMetrictHandler(reqName string) *HTTPMetrictHandler {
	return &HTTPMetrictHandler{
		reqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_total", reqName),
			Name:      "requests",
		}),
		errCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_total", reqName),
			Name:      "errors",
		}),
		reqLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: fmt.Sprintf("http_%s_latency", reqName),
			Name:      "requests",
			Buckets:   []float64{0.1, 0.5, 1},
		}),
	}
}

func (h *HTTPMetrictHandler) instrument(next HTTPFunc) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()

			logrus.WithFields(logrus.Fields{
				"method":  r.Method,
				"path":    r.URL.Path,
				"latency": fmt.Sprintf("%.2fs", latency),
				"error":   err,
			}).Info("request processed")

			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()

			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())

		err = next(w, r)
		return err
	}
}

func handleGetInvoice(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return APIError{Code: http.StatusMethodNotAllowed, Err: fmt.Errorf("method not supported (%s)", http.MethodGet)}
		}

		q, ok := r.URL.Query()["obu"]
		if !ok {
			return APIError{Code: http.StatusBadRequest, Err: fmt.Errorf("missing OBU ID")}
		}

		OBUID, err := strconv.Atoi(q[0])
		if err != nil {
			return APIError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid OBU ID (%s)", q[0])}
		}

		invoice, err := svc.CalculateInvoice(OBUID)
		if err != nil {
			return APIError{Code: http.StatusInternalServerError, Err: err}
		}

		return writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return APIError{Code: http.StatusMethodNotAllowed, Err: fmt.Errorf("method not supported (%s)", r.Method)}
		}

		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return APIError{Code: http.StatusBadRequest, Err: err}
		}
		if err := svc.AggregateDistance(distance); err != nil {
			return APIError{Code: http.StatusInternalServerError, Err: err}
		}

		return writeJSON(w, http.StatusCreated, nil)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
