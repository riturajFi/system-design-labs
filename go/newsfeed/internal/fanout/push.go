package fanout

import (
	"context"
	"encoding/json"

	"newsfeed/internal/graph"
	"newsfeed/internal/models"
	"newsfeed/internal/queue"
)

// PushStrategy implements fanout-on-write.
type PushStrategy struct {
	graph graph.Store
	queue queue.Queue
}

type FanoutMessage struct {
	UserID string
	PostID string
}

func NewPushStrategy(
	graph graph.Store,
	queue queue.Queue,
) *PushStrategy {
	return &PushStrategy{
		graph: graph,
		queue: queue,
	}
}

func (p *PushStrategy) Distribute(ctx context.Context, post models.Post) error {
	friends := p.graph.Friends(ctx, post.AuthorID)

	for _, friendID := range friends {
		payload, err := json.Marshal(FanoutMessage{
			UserID: friendID,
			PostID: post.ID,
		})
		if err != nil {
			return err
		}

		msg := queue.Message{
			Key:   friendID,
			Value: string(payload),
		}
		if err := p.queue.Publish(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}
