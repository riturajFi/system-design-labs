package auth

import "context"

type MockAuthenticator struct {}

func (m *MockAuthenticator) Authenticate (
	_ context.Context,
	token string,
) (string, bool) {
	const prefix = "token-"

	if(len(token) <= len(prefix)){
		return "", false
	}

	if(token[:len(prefix)] != prefix){

		return "", false
	}

	return token[len(prefix):], true
}