package aggservice

import (
	"fmt"
	"log"
	"os"
	"tolling/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func makeStore() Storer {
	storeType := os.Getenv("STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("store type not supported - %s", storeType)
		return nil
	}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemoryStore) Get(ID int) (float64, error) {
	dist, ok := m.data[ID]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for obu id [%d]", ID)
	}
	return dist, nil
}
