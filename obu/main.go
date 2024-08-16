package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
	"toll-calculator/types"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func genCoord() float64 {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	n := float64(rand.Intn(100) + 1)
	f := r.Float64()
	return n + f
}

func genLocation() (float64, float64) {
	return genCoord(), genCoord()
}

func generateOBUIDS(n int) []int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	ids := make([]int, 0, n)

	for i := 0; i < n; i++ {
		ids = append(ids, r.Intn(math.MaxInt))
	}

	return ids
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	wsEndpoint := fmt.Sprintf("ws://127.0.0.1%s/ws", os.Getenv("DATA_RECEIVER_WS_PORT"))
	OBUIDS := generateOBUIDS(20)

	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)

	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 0; i < len(OBUIDS); i++ {
			lat, long := genLocation()
			data := types.OBUData{
				OBUID: OBUIDS[i],
				Lat:   lat,
				Long:  long,
			}

			if err := conn.WriteJSON(&data); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%+v\n", data)
		}

		time.Sleep(time.Second * 5)
	}
}
