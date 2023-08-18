package interfaces

import "time"

type QueueInterface interface {
	Enqueue(value string, timestamp time.Time)
	Dequeue() interface{}
	IsEmpty() bool
	Size() int
	First() *string
	ReadAll() []string
}
