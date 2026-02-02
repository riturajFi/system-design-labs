package model

type Notification struct {
	EventID    string
	UserID     string
	Channel    Channel
	TemplateID string
	Params     map[string]string
	CreatedAt  int64
}
