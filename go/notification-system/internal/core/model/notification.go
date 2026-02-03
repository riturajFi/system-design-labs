package model

type Notification struct {
	EventID    string
	UserID     int64
	Channel    Channel
	TemplateID string
	Params     map[string]string
	CreatedAt  int64
}
