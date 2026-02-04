package queue

import "context"

type Message struct {
	Key string
	Value string
}

type Queue interface {
	Publish(ctx context.Context, msg Message) error
	Consume(ctx context.Context) (Message, error)
}