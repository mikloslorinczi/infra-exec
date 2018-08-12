package syncmap

import (
	"sync"
)

// MapI interface requires structures implementing it to have a
// GetMap function which returns the map
// SetKey function which sets the given key to the given value
// GetKey function which returns the value of the given key
// DeleteKey function which renders the given key's value to its initial zero (empty string).
type MapI interface {
	GetMap() map[string]string
	SetKey(key string, value string)
	GetKey(key string) string
	DeleteKey(key string)
}

// SafeMap implements the MapI interface. GetMap, SetKey, GetKey and DeleteKey
// can be callen on its string-map called data. SafeMap uses sync's RWMutex
// to manage safe concurent calls.
type SafeMap struct {
	rwMutex sync.RWMutex
	data    map[string]string
}

// NewSafeMap returns a pointer to a SafeMap
func NewSafeMap() MapI {
	return &SafeMap{data: make(map[string]string)}
}

// GetMap returns the map using Read Lock.
func (m *SafeMap) GetMap() map[string]string {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	return m.data
}

// SetKey sets a given key to a given value using Lock.
func (m *SafeMap) SetKey(key string, value string) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.data[key] = value
}

// GetKey returns a given key's value using Read Lock.
func (m *SafeMap) GetKey(key string) string {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	return m.data[key]
}

// DeleteKey renders a given key's value to its initial zero (empty string)
func (m *SafeMap) DeleteKey(key string) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	delete(m.data, key)
}
