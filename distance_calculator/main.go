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
	httpAggClient := client.NewHTTPClient("http://localhost:5000")
	grpcAggClient, err := client.NewGRPCClient("localhost:5001")

	if err != nil {
		log.Fatal(err)
	}

	_ = grpcAggClient
	_ = httpAggClient
	c, err := NewKafkaConsumer(kafkaTopic, srv, httpAggClient)

	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
