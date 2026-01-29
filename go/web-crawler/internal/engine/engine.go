package engine

import (
	"fmt"
	"web-crawler/internal/fetcher"
	"web-crawler/internal/model"
)

type Engine struct {
	fetcher fetcher.Fetcher
}

func New(fetcher fetcher.Fetcher) *Engine {
	return &Engine{
		fetcher: fetcher,
	}
}

func (e *Engine) Run(seedUrl string) error {

	req := model.CrawlRequest{
		URL: seedUrl,
	}

	result, err := e.fetcher.Fetch(req)
	if  err != nil  {
		return err
	}

	fmt.Printf("URL: %s\n", result.URL)
	fmt.Printf("Status: %d\n", result.StatusCode)
	fmt.Printf("Bytes: %d\n", len(result.Body))

	return  nil
}