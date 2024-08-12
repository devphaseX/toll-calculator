package main

import (
	"log"
	"toll-calculator/aggregator/client"
)

const kafkaTopic = "obudata"

func main() {

	var (
		srv CalculatorServicer
	)
	srv = NewCalculatorService()
	srv = NewLogMiddleware(srv)
	aggClient := client.NewClient("http://localhost:5000/agg")
	c, err := NewKafkaConsumer(kafkaTopic, srv, aggClient)

	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
