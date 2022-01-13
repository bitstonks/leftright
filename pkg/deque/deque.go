package deque

import (
	"errors"
)

// Operation is the generic input accepted by the left-right data structure. Defined as an alias for readability.
type Operation interface{}

// Deque is a generic double ended queue implemented with a ring buffer.
type Deque struct {
	// data is the underlying data storage. Starts at a given capacity and is doubled whenever we fill it.
	data []Operation
	// head is the index of the first valid element in data, if any.
	head int
	// tail is the index of the last valid element in data, if any.
	tail int
	// count keeps track of the number of elements we're storing in data.
	count int
}

// NewDequeWithCapacity creates a new deque with an initial capacity of `cap`.
func NewDequeWithCapacity(cap int) Deque {
	return Deque{
		data:  make([]Operation, cap),
		head:  0,
		tail:  0,
		count: 0,
	}
}

// NewDeque creates a new Deque with initial capacity 16.
// Why 16? https://youtu.be/0obMRztklqU
func NewDeque() Deque {
	return NewDequeWithCapacity(16)
}

// Len returns the number of elements in the deque.
func (q *Deque) Len() int {
	return q.count
}

// checkCapacity doubles the capacity of deque if it ran out of space.
func (q *Deque) checkCapacity() {
	if q.count < len(q.data) {
		return
	}

	if len(q.data) == 0 {
		q.data = make([]Operation, 1)
		return
	}

	if q.head == 0 {
		// If head is at zero, we don't have to do any work, just append empty slots.
		q.data = append(q.data, make([]Operation, q.count)...)
		return
	}
	// Since we're using a ring buffer we have to move elements around to ensure a sequential buffer.
	// Because data is full we know that it looks something like: 6 7 T H 2 3 4 5
	// but we want to change it to: . . . H 2 3 4 5 6 7 T . . . . .
	// We do this by duplicating it and then setting unneeded elements to nil to prevent memory leaks.
	q.data = append(q.data, q.data...)
	q.tail += q.count
	for i := 0; i < q.head; i++ {
		q.data[i] = nil
	}
	for i := q.tail+1; i < len(q.data); i++ {
		q.data[i] = nil
	}
}

// PushBack inserts a new element at the end.
func (q *Deque) PushBack(op Operation) {
	q.checkCapacity()
	if q.count == 0 {
		q.head = 0
		q.tail = 0
	} else {
		q.tail += 1
		q.tail %= len(q.data)
	}
	q.data[q.tail] = op
	q.count += 1
}

// PushFront inserts a new element at the start.
func (q *Deque) PushFront(op Operation) {
	q.checkCapacity()
	if q.count == 0 {
		q.head = 0
		q.tail = 0
	} else {
		q.head += len(q.data) - 1
		q.head %= len(q.data)
	}
	q.data[q.head] = op
	q.count += 1
}

// PopBack removes and returns the last element in deque or returns an error if deque is empty.
func (q *Deque) PopBack() (Operation, error) {
	if q.count == 0 {
		return nil, errors.New("cannot PopBack because deque is empty")
	}
	val := q.data[q.tail]
	q.data[q.tail] = nil
	q.count -= 1
	q.tail += len(q.data) - 1
	q.tail %= len(q.data)
	return val, nil
}

// PopFront removes and returns the first element in deque or returns an error if deque is empty.
func (q *Deque) PopFront() (Operation, error) {
	if q.count == 0 {
		return nil, errors.New("cannot PopFront because deque is empty")
	}
	val := q.data[q.head]
	q.data[q.head] = nil
	q.count -= 1
	q.head += 1
	q.head %= len(q.data)
	return val, nil
}
