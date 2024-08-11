package main

import (
	"time"
	"toll-calculator/types"

	"github.com/sirupsen/logrus"
)

type LogginMiddleware struct {
	next DataProducer
}

func NewLogginMiddleware(next DataProducer) *LogginMiddleware {
	return &LogginMiddleware{
		next: next,
	}
}

func (l *LogginMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"OBUID": data.OBUID,
			"Lat":   data.Lat,
			"Long":  data.Long,
			"took":  time.Since(start),
		}).Info("producing to kafka")
	}(time.Now())

	return l.next.ProduceData(data)
}
