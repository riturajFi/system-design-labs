package model


type CrawlRequest struct {
	URL string
}

type FetchResult struct {
	URL string
	StatusCode int
	Body []byte
}