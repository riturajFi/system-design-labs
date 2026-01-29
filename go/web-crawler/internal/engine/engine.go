package engine

import (
	"fmt"
	"web-crawler/internal/dedupe"
	"web-crawler/internal/fetcher"
	"web-crawler/internal/frontier"
	"web-crawler/internal/parser"
)

type Engine struct {
	fetcher  fetcher.Fetcher
	frontier frontier.Frontier
	deduper  dedupe.Deduper
	parser   parser.Parser
}

func New(
	fetcher fetcher.Fetcher,
	frontier frontier.Frontier,
	deduper dedupe.Deduper,
	parser parser.Parser,
) *Engine {
	return &Engine{
		fetcher:  fetcher,
		frontier: frontier,
		deduper:  deduper,
		parser:   parser,
	}
}

func (e *Engine) Run() error {
	for {
		req, ok := e.frontier.Pop()
		if !ok {
			return nil
		}

		if e.deduper.Seen(req) {
			continue
		}

		e.deduper.Mark(req)

		result, err := e.fetcher.Fetch(req)
		if err != nil {
			fmt.Println("fetch error:", err)
			continue
		}

		fmt.Printf("Fetched: %s (%d bytes)\n", result.URL, len(result.Body))

		children, err := e.parser.Parse(result.URL, result.Body)
		if err != nil {
			fmt.Println("parse error:", err)
			continue
		}

		for _, child := range children {
			e.frontier.Push(child)
		}
	}
}
