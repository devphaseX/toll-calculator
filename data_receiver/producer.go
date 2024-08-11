package main

import (
	"encoding/json"
	"fmt"
	"toll-calculator/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	ProduceData(types.OBUData) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	p.Flush(0)
	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (k *KafkaProducer) ProduceData(data types.OBUData) error {
	// Produce messages to topic (asynchronously)
	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
}
