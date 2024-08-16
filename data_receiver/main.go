package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"toll-calculator/types"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const kafkaTopic = "obudata"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	recv, err := NewDataReceiver()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(os.Getenv("DATA_RECEIVER_WS_PORT"), nil)
}

type DataReceiver struct {
	conn     *websocket.Conn
	producer DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := NewKafkaProducer(kafkaTopic)

	if err != nil {
		return nil, err
	}

	p = NewLogginMiddleware(p)

	return &DataReceiver{
		producer: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.producer.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	if dr.conn == nil {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Fatal(err)
		}

		dr.conn = conn
	}

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU client connected")
	for {
		var data types.OBUData

		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Printf("read error: %v\n", err)
			continue
		}

		if err := dr.produceData(data); err != nil {
			log.Printf("kafka produce error: %v", err)
			continue
		}
	}
}
