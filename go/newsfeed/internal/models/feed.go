package models

import "time"

// FeedEntry represents ONE post in ONE user's feed.
// IDs only. No embedded objects.
type FeedEntry struct {
	UserID    string
	PostID    string
	CreatedAt time.Time
}
