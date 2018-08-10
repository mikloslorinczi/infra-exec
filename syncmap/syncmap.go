package syncmap

import (
	"sync"
)

// HashMapI interface exposes getMap and addToMap methods for the client,
// in order to interact with the underlying map.
type HashMapI interface {
	GetMap() map[string]string
	AddToMap(key string, value string)
}

type hashMap struct {
	rwMutex sync.RWMutex
	data    map[string]string
}

// NewHashMap returns a pointer to a hashMap
func NewHashMap() HashMapI {
	return &hashMap{data: make(map[string]string)}
}

func (h *hashMap) GetMap() map[string]string {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()
	return h.data
}

func (h *hashMap) AddToMap(key string, value string) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()
	h.data[key] = value
}
