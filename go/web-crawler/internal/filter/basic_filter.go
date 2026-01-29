package filter

import (
	"net/url"
	"web-crawler/internal/model"
)

type BasicFilter struct{}

func NewBasicFilter() *BasicFilter {
	return &BasicFilter{}
}

func (f *BasicFilter) Allow(req model.CrawlRequest) bool {
	u, err := url.Parse(req.URL)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}
