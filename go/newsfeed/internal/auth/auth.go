package auth

import "context"

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (userID string, ok bool)
}
