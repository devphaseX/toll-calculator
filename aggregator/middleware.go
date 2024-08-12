package main

import (
	"time"
	"toll-calculator/types"

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
			"took": time.Since(start),
			"err":  err,
		}).Info("aggregating distance")
	}(time.Now())

	err = m.next.AggregateDistance(distance)
	return err
}
