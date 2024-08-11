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
}

func NewKafkaProducer() (DataProducer, error) {
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
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	p.Flush(0)
	return &KafkaProducer{
		producer: p,
	}, nil
}

func (k *KafkaProducer) ProduceData(data types.OBUData) error {
	// Produce messages to topic (asynchronously)
	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	topic := kafkaTopic
	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
}
