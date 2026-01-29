package robots

import "web-crawler/internal/model"

// Policy decides if a crawl request is allowed by robots.txt.
type Policy interface {
	Allowed(req model.CrawlRequest) bool
}
