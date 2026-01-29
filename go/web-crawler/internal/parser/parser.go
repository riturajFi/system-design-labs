package parser

import "web-crawler/internal/model"

// Parser extracts URLs from fetched content.
type Parser interface {
	Parse(baseURL string, body []byte) ([]model.CrawlRequest, error)
}
