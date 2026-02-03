package mock

import (
	"context"
	"errors"

	"notification-system/internal/core/model"
)

type Provider struct {
	Fail bool
}

func New(fail bool) *Provider {
	return &Provider{Fail: fail}
}

func (p *Provider) Send(
	_ context.Context,
	n model.Notification,
) error {

	if p.Fail {
		return errors.New("email provider simulated failure")
	}
	return nil
}
