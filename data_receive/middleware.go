package main

import (
	"time"
	"tolling/types"

	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	next DataProducer
}

func NewLoggerMiddleware(next DataProducer) *LoggerMiddleware {
	return &LoggerMiddleware{
		next: next,
	}
}

func (lm *LoggerMiddleware) ProducerData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obu_id": data.OBUID,
			"lat":    data.Lat,
			"long":   data.Long,
			"took":   time.Since(start),
		}).Info("Received data")
	}(time.Now())

	return lm.next.ProducerData(data)
}
