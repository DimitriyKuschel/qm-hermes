package queues

import (
	"sync"
	"time"
)

type Item struct {
	message   string
	timestamp time.Time
	next      *Item
}

type FIFOQueue struct {
	first *Item
	last  *Item
	mu    sync.Mutex
}

func (q *FIFOQueue) Enqueue(value string, timestamp time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()

	newItem := &Item{message: value, timestamp: timestamp}

	if q.first == nil {
		// FIFOQueue is empty
		q.first = newItem
		q.last = newItem
	} else {
		if q.first.timestamp.After(newItem.timestamp) {
			// new item is earlier than the first item
			newItem.next = q.first
			q.first = newItem
		} else {
			// find the first item that has a timestamp later than the new item
			current := q.first
			for current.next != nil && current.next.timestamp.Before(newItem.timestamp) {
				current = current.next
			}

			// insert the new item
			newItem.next = current.next
			current.next = newItem

			if newItem.next == nil {
				// update last item
				q.last = newItem
			}
		}
	}
}

func (q *FIFOQueue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.first == nil {
		// FIFOQueue is empty
		return nil
	}

	item := q.first
	q.first = item.next

	if q.first == nil {
		// FIFOQueue is now empty
		q.last = nil
	}

	return item.message
}

func (q *FIFOQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.first == nil
}

func (q *FIFOQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	size := 0
	item := q.first
	for item != nil {
		size++
		item = item.next
	}

	return size
}

func (q *FIFOQueue) First() *string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.first == nil {
		// FIFOQueue is empty
		return nil
	}

	return &q.first.message
}

func (q *FIFOQueue) ReadAll() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	response := []string{}
	if q.first == nil {
		// FIFOQueue is empty
		return response
	}
	item := q.first

	for item != nil {
		response = append(response, item.message)
		item = item.next
	}

	return response
}
