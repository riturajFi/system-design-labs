package queue

import (
	"context"
)

type MemoryQueue struct {
	ch chan Message
}

func NewMemoryQueue(buffer int) *MemoryQueue {
	return &MemoryQueue{
		ch: make(chan Message, buffer),
	}
}

func (q *MemoryQueue) Publish (ctx context.Context, msg Message) error {
	select {
	case q.ch <-msg:
		return  nil
	case <- ctx.Done():
		return ctx.Err()
	}
}

func (q *MemoryQueue) Consume(ctx context.Context) (Message, error) {
	select {
	case msg := <- q.ch:
		return msg, nil
	case <- ctx.Done():
		return Message{}, ctx.Err()
	}
}