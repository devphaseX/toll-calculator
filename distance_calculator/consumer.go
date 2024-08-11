package main

import (
	"encoding/json"
	"toll-calculator/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type kafkaConsumer struct {
	consumer   *kafka.Consumer
	IsRunnning bool
	srv        CalculatorServicer
}

func NewKafkaConsumer(topic string, srv CalculatorServicer) (*kafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		consumer: c,
		srv:      srv,
	}, nil
}

func (c *kafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.IsRunnning = true
	c.readMessageLoop()
}

func (c *kafkaConsumer) readMessageLoop() {

	for c.IsRunnning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consumer error: %s\n", err.Error())
			continue
		}

		var data types.OBUData

		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s\n", err.Error())
			continue
		}

		distance, err := c.srv.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculator error:  %s\n", err)
			continue
		}

		_ = distance
	}
}
