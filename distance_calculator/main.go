package main

import "log"

const kafkaTopic = "obudata"

func main() {

	var (
		srv CalculatorServicer
	)
	srv = NewCalculatorService()
	srv = NewLogMiddleware(srv)
	c, err := NewKafkaConsumer(kafkaTopic, srv)

	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
