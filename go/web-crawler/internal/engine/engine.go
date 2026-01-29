package engine

import (
	"fmt"
	"web-crawler/internal/dedupe"
	"web-crawler/internal/fetcher"
	"web-crawler/internal/frontier"
)

type Engine struct {
	fetcher  fetcher.Fetcher
	frontier frontier.Frontier
	deduper  dedupe.Deduper
}

func New(fetcher fetcher.Fetcher, frontier frontier.Frontier, deduper dedupe.Deduper) *Engine {
	return &Engine{
		fetcher:  fetcher,
		frontier: frontier,
		deduper:  deduper,
	}
}

func (e *Engine) Run() error {
	for {
		req, ok := e.frontier.Pop()
		if !ok {
			return nil
		}

		if e.deduper.Seen(req) {
			fmt.Println("skipping already seen URL:", req.URL)
			continue
		}

		e.deduper.Mark(req)

		result, err := e.fetcher.Fetch(req)
		if err != nil {
			fmt.Println("fetch error:", err)
			continue
		}

		fmt.Printf("URL: %s\n", result.URL)
		fmt.Printf("Status: %d\n", result.StatusCode)
		fmt.Printf("Bytes: %d\n\n", len(result.Body))
	}
}
