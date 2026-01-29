package engine

import (
	"fmt"
	"sync"
	"web-crawler/internal/dedupe"
	"web-crawler/internal/fetcher"
	"web-crawler/internal/filter"
	"web-crawler/internal/frontier"
	"web-crawler/internal/model"
	"web-crawler/internal/parser"
)

const workerCount = 4

type Engine struct {
	fetcher  fetcher.Fetcher
	frontier frontier.Frontier
	deduper  dedupe.Deduper
	parser   parser.Parser
	filter   filter.Filter
}

func New(
	fetcher fetcher.Fetcher,
	frontier frontier.Frontier,
	deduper dedupe.Deduper,
	parser parser.Parser,
	flt filter.Filter,
) *Engine {
	return &Engine{
		fetcher:  fetcher,
		frontier: frontier,
		deduper:  deduper,
		parser:   parser,
		filter:   flt,
	}
}

func (e *Engine) Run() error {
	workCh := make(chan model.CrawlRequest)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			e.worker(id, workCh)
		}(i)
	}

	// Dispatcher loop (single-threaded)
	for {
		req, ok := e.frontier.Pop()
		if !ok {
			close(workCh)
			break
		}

		if e.deduper.Seen(req) {
			continue
		}

		e.deduper.Mark(req)
		workCh <- req
	}

	wg.Wait()
	return nil
}

func (e *Engine) worker(id int, workCh <-chan model.CrawlRequest) {
	for req := range workCh {
		result, err := e.fetcher.Fetch(req)
		if err != nil {
			fmt.Printf("[worker %d] fetch error: %v\n", id, err)
			// Always release the host slot on failure.
			e.frontier.Done(req)
			continue
		}

		fmt.Printf("[worker %d] fetched %s (%d bytes)\n", id, result.URL, len(result.Body))

		children, err := e.parser.Parse(result.URL, result.Body)
		if err != nil {
			fmt.Printf("[worker %d] parse error: %v\n", id, err)
			// Parsing failure still completes the request lifecycle.
			e.frontier.Done(req)
			continue
		}

		for _, child := range children {
			if !e.filter.Allow(child) {
				continue
			}
			e.frontier.Push(child)
		}

		// Notify frontier that this request is fully processed.
		e.frontier.Done(req)
	}
}
