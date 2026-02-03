package memory

import (
	"context"
	"errors"

	"notification-system/internal/core/model"
)

type Store struct {
	users   map[int64]model.User
	devices map[int64][]model.Device
}

func New(
	users []model.User,
	devices []model.Device,
) *Store {
	u := make(map[int64]model.User)
	for _, user := range users {
		u[user.ID] = user
	}

	d := make(map[int64][]model.Device)
	for _, device := range devices {
		d[device.UserID] = append(d[device.UserID], device)
	}

	return &Store{
		users:   u,
		devices: d,
	}
}

func (s *Store) GetUser(
	_ context.Context,
	userID int64,
) (model.User, error) {
	user, ok := s.users[userID]
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *Store) GetDevicesByUser(
	_ context.Context,
	userID int64,
) ([]model.Device, error) {
	return s.devices[userID], nil
}
