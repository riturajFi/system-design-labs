package dedupe

import "web-crawler/internal/model"

// Deduper decides whether a URL has been seen before.
type Deduper interface {
	Seen(req model.CrawlRequest) bool
	Mark(req model.CrawlRequest)
}
