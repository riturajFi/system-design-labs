package model

// CrawlRequest represents a single URL crawl task with priority metadata.
type CrawlRequest struct {
	URL      string
	Priority Priority
}

type FetchResult struct {
	URL        string
	StatusCode int
	Body       []byte
}
