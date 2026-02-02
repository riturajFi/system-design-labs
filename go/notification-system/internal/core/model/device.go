package model

type Device struct {
	ID       int64
	UserID   int64
	Token    string
	Platform Channel
	LastSeen int64
}
