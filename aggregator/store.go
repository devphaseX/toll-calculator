package main

import (
	"fmt"
	"toll-calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(data types.Distance) error {
	m.data[data.OBUID] += data.Value
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	dist, ok := m.data[id]

	if !ok {
		return 0.0, fmt.Errorf("could not find distance with obu id %d", id)
	}

	return dist, nil
}
