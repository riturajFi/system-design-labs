package model

type EventType string

const (
	EventPending   EventType = "pending"
	EventSent      EventType = "sent"
	EventDelivered EventType = "delivered"
	EventError     EventType = "error"
	EventClick     EventType = "click"
	EventUnsub     EventType = "unsubscribe"
)

type NotificationEvent struct {
	EventID string
	UserID  string
	Channel Channel
	Type    EventType
	Message string
	AtUnix  int64
}
