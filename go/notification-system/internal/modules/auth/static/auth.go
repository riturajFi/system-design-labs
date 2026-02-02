package static

import (
	"context"
	"errors"
)

type StaticAuthenticator struct {
	keys map[string]string
}

func New(keys map[string]string) *StaticAuthenticator {
	return &StaticAuthenticator{keys: keys}
}

func (a *StaticAuthenticator) Authenticate(
	_ context.Context,
	appKey string,
	appSecret string,
) error {
	secret, ok := a.keys[appKey]
	if !ok {
		return errors.New("invalid app key")
	}

	if secret != appSecret {
		return errors.New("invalid app secret")
	}

	return nil
}
