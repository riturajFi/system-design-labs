package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type Rendered struct {
	Channel model.Channel
	Title   string
	Body    string
}

type TemplateEngine interface {
	Render(
		ctx context.Context,
		templateID string,
		channel model.Channel,
		params map[string]string,
	) (Rendered, error)
}
