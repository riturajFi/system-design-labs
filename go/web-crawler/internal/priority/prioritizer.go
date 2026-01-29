package priority

import "web-crawler/internal/model"

// Prioritizer assigns a scheduling priority to a crawl request.
type Prioritizer interface {
	Assign(req model.CrawlRequest) model.CrawlRequest
}
