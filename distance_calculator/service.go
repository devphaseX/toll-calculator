package main

import (
	"math"
	"toll-calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct {
	prev []float64
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (c *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if c.prev != nil {
		lat, long := c.prev[0], c.prev[1]
		distance = calculateDistance(lat, long, data.Lat, data.Long)
	}

	c.prev = []float64{data.Lat, data.Long}

	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
