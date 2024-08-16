package main

import (
	"fmt"
	"time"
	"toll-calculator/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {

	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)

		if err == nil {
			distance = invoice.TotalDistance
			amount = invoice.TotalAmount
		}

		logrus.WithFields(logrus.Fields{
			"took":          time.Since(start),
			"err":           err,
			"obuID":         obuID,
			"TotalDistance": distance,
			"TotalAmount":   amount,
		}).Info("CalculateInvoice")
	}(time.Now())

	invoice, err = m.next.CalculateInvoice(obuID)
	return invoice, err
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"obuID": distance.OBUID,
			"err":   err,
		}).Info("aggregating distance")
	}(time.Now())

	err = m.next.AggregateDistance(distance)
	return err
}

type MetricMiddleware struct {
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           Aggregator
}

func NewMetricMiddleware(next Aggregator) *MetricMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "error_request_counter",
		Name:      "aggregator",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "error_request_counter",
		Name:      "calculator",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculator",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregator",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculator",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &MetricMiddleware{
		next:           next,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqCounterCalc: reqCounterCalc,
		reqCounterAgg:  reqCounterAgg,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}

func (m *MetricMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		fmt.Println("prometheus")
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()

		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(obuID)
	return
}

func (m *MetricMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()

		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return err
}
