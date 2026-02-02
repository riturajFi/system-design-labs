package ports

import "context"

type Authenticator interface {
	Authenticate(ctx context.Context, appKey, appSecret string) error
}
