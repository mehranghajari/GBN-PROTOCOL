// Package queue creates a ItemQueue data structure for the Item type
package queue

// Item the type of the queue

// ItemQueue the queue of Items
type Windows struct {
	items []string
	size int
	elementsNumber int
}

// New creates a new ItemQueue
func (s *Windows) New(size int) *Windows {
	s.size = size
	s.elementsNumber = 0
	s.items = make([]string, size)
	return s
}

// Enqueue adds an Item to the end of the queue
func (s *Windows) Enqueue(t []byte) {
	if !s.IsFull() {
		s.items = append(s.items, string(t))
		s.elementsNumber++
	}

}

// Dequeue removes an Item from the start of the queue
func (s *Windows) Dequeue() []byte {
	if !s.IsEmpty() {
		item := s.items[0]
		s.items = s.items[1:len(s.items)]
		s.elementsNumber--
		return []byte(item)
	}
	return nil
}

// Front returns the item next in the queue, without removing it
func (s *Windows) Front() []byte {
	item := s.items[0]
	return []byte(item)
}

// IsEmpty returns true if the queue is empty
func (s *Windows) IsEmpty() bool {
	return s.elementsNumber == 0
}

// Size returns the number of Items in the queue
func (s *Windows) Size() int {
	return s.elementsNumber
}
func (s *Windows) IsFull() bool{
	return s.elementsNumber == s.size

}
func (s *Windows) GetArray() []string {
	return s.items
}