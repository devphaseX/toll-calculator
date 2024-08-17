package main

import (
	"fmt"
	"log"
	"os"
	"toll-calculator/aggregator/client"

	"github.com/joho/godotenv"
)

const kafkaTopic = "obudata"

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatal(err)
	}

	var (
		srv CalculatorServicer
	)
	srv = NewCalculatorService()
	srv = NewLogMiddleware(srv)
	httpAggClient := client.NewHTTPClient(fmt.Sprintf("http://localhost%s", os.Getenv("AGG_HTTP_PORT")))
	grpcAggClient, err := client.NewGRPCClient(fmt.Sprintf("localhost%s", os.Getenv("AGG_GRPC_PORT")))

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
