package main

import (
	"time"
	"toll-calculator/types"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) *LogMiddleware {
	return &LogMiddleware{
		next,
	}
}

func (l *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"tooks":    time.Since(start),
			"distance": dist,
			"err":      err,
		}).Info("calculate distance")
	}(time.Now())

	dist, err = l.next.CalculateDistance(data)
	return
}
