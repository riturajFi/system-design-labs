package models

import "time"

type Post struct {
	ID        string
	AuthorID  string
	Content   string
	CreatedAt time.Time
}
