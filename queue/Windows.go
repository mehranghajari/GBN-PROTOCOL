// Package queue creates a ItemQueue data structure for the Item type
package queue

import (
	"sync"
)

// Item the type of the queue

// ItemQueue the queue of Items
type Windows struct {
	items [][]byte
	lock  sync.RWMutex
	size int
}

// New creates a new ItemQueue
func (s *Windows) New(size int) *Windows {
	s.size = size
	s.items = make([][]byte, s.size)
	return s
}

// Enqueue adds an Item to the end of the queue
func (s *Windows) Enqueue(t []byte) {
	s.lock.Lock()
	s.items=append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Item from the start of the queue
func (s *Windows) Dequeue() []byte {
	s.lock.Lock()
	if !s.IsEmpty() {
		item := s.items[0]
		s.items = s.items[1:len(s.items)]
		s.lock.Unlock()
		return item
	}
	return nil
}

// Front returns the item next in the queue, without removing it
func (s *Windows) Front() []byte {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return item
}

// IsEmpty returns true if the queue is empty
func (s *Windows) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of Items in the queue
func (s *Windows) Size() int {
	return len(s.items)
}