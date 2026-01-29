package filter

import "web-crawler/internal/model"

// Filter decides whether a URL should be crawled.
type Filter interface {
	Allow(req model.CrawlRequest) bool
}
