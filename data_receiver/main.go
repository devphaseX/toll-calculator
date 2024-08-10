package main

import (
	"fmt"
	"log"
	"net/http"
	"toll-calculator/types"

	"github.com/gorilla/websocket"
)

func main() {
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":3000", nil)
}

type DataReceiver struct {
	conn *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{}
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

		fmt.Printf("receive OBU data from [%d] :: <lat %.2f, long %2.f>\n", data.OBUID, data.Lat, data.Long)
	}
}
