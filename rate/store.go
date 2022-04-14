package rate

import (
	"sync"
	"time"
)

// Store provides a method to store counters.
type Store interface {
	Increment(key string, expiry int64) (int, error)
}

// memoryRecord represents a record in memory.
type memoryRecord struct {
	value  int
	expiry int64
}

// newMemoryRecord returns an instance of memoryRecord.
func newMemoryRecord(value int, expiry int64) *memoryRecord {
	return &memoryRecord{
		value:  value,
		expiry: expiry,
	}
}

// MemoryStore stores counters in memory and implements Store interface to
// increment them.
type MemoryStore struct {
	store map[string]*memoryRecord
	mutex sync.Mutex
}

// NewMemoryStore returns an instance of MemoryStore. MemoryStore instances
// created using this function also provides eviction functionality which
// removes the old items in memory based on their expiry.
func NewMemoryStore(evictionFrequency time.Duration) *MemoryStore {
	s := &MemoryStore{
		store: map[string]*memoryRecord{},
	}

	go func() {
		for now := range time.Tick(evictionFrequency) {
			s.mutex.Lock()
			for key, record := range s.store {
				if now.Unix() > record.expiry {
					delete(s.store, key)
				}
			}
			s.mutex.Unlock()
		}
	}()

	return s
}

// Increment the counter for the given key and set its expiry.
func (s *MemoryStore) Increment(key string, expiry int64) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := s.store[key]
	if r == nil {
		r = newMemoryRecord(0, expiry)
		s.store[key] = r
	}
	r.value += 1

	return r.value, nil
}
