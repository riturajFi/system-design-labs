package fetcher

import (
	"web-crawler/internal/model"
)

type Fetcher interface {
	Fetch(req model.CrawlRequest) (*model.FetchResult, error)
}