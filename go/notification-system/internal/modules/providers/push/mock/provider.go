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

	switch n.Channel {
	case model.ChannelPushIOS, model.ChannelPushAndroid:
		if p.Fail {
			return errors.New("push provider simulated failure")
		}
		return nil
	default:
		return errors.New("unsupported push channel")
	}
}
