package basic

import (
	"context"
	"errors"

	"notification-system/internal/core/model"
	"notification-system/internal/core/ports"
)

type Resolver struct {
	contacts ports.ContactStore
}

func New(contacts ports.ContactStore) *Resolver {
	return &Resolver{contacts: contacts}
}

func (r *Resolver) Resolve(
	ctx context.Context,
	n model.Notification,
) ([]ports.Target, error) {
	user, err := r.contacts.GetUser(ctx, n.UserID)
	if err != nil {
		return nil, err
	}

	switch n.Channel {
	case model.ChannelEmail:
		if user.Email == "" {
			return nil, errors.New("user has no email")
		}
		return []ports.Target{
			{
				Channel: model.ChannelEmail,
				Address: user.Email,
				User:    &user,
			},
		}, nil

	case model.ChannelSMS:
		if user.PhoneNumber == "" {
			return nil, errors.New("user has no phone number")
		}
		return []ports.Target{
			{
				Channel: model.ChannelSMS,
				Address: user.PhoneNumber,
				User:    &user,
			},
		}, nil

	case model.ChannelPushIOS, model.ChannelPushAndroid:
		devices, err := r.contacts.GetDevicesByUser(ctx, n.UserID)
		if err != nil {
			return nil, err
		}

		var targets []ports.Target
		for _, d := range devices {
			if d.Platform == n.Channel {
				device := d
				targets = append(targets, ports.Target{
					Channel: n.Channel,
					Address: d.Token,
					Device:  &device,
					User:    &user,
				})
			}
		}

		if len(targets) == 0 {
			return nil, errors.New("no devices for channel")
		}
		return targets, nil

	default:
		return nil, errors.New("unknown channel")
	}
}
