package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const sentInterval = 60

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

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

	OBUIDS := generateOBUIDS(20)

	for {
		for i := 0; i < len(OBUIDS); i++ {
			lat, long := genLocation()
			data := OBUData{
				OBUID: OBUIDS[i],
				Lat:   lat,
				Long:  long,
			}

			fmt.Printf("%+v\n", data)
		}

		time.Sleep(sentInterval)
	}
}
