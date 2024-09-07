package main

import (
	"time"
	"tolling/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	errCounterAgg prometheus.Counter
	errCounterCal prometheus.Counter
	reqCounterAgg prometheus.Counter
	reqCounterCal prometheus.Counter
	reqLatencyAgg prometheus.Histogram
	reqLatencyCal prometheus.Histogram
	next          Aggregator
}

// NewMetrictsMeddleware push the metrics to prometheus
func NewMetrictsMeddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregator",
	})

	errCounterCal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})

	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregator",
	})

	reqCounterCal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	reqLatencyCal := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &MetricsMiddleware{
		errCounterAgg: errCounterAgg,
		errCounterCal: errCounterCal,
		reqCounterAgg: reqCounterAgg,
		reqCounterCal: reqCounterCal,
		reqLatencyAgg: reqLatencyAgg,
		reqLatencyCal: reqLatencyCal,
		next:          next,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *MetricsMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCal.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterCal.Inc()
		if err != nil {
			m.errCounterCal.Inc()
		}
	}(time.Now())
	_, err = m.next.CalculateInvoice(obuId)
	return
}

// LogMiddleware middleware for log aggregator
type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)

		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}

		logrus.WithFields(logrus.Fields{
			"took":           time.Since(start),
			"err":            err,
			"obu_id":         obuId,
			"total_distance": distance,
			"total_amount":   amount,
		}).Info("CalculateInvoice")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuId)
	return
}
