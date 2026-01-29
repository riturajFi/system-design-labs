package main

import (
	"fmt"
	"os"
	"web-crawler/internal/dedupe"
	"web-crawler/internal/engine"
	"web-crawler/internal/fetcher"
	"web-crawler/internal/filter"
	"web-crawler/internal/frontier"
	"web-crawler/internal/model"
	"web-crawler/internal/parser"
	"web-crawler/internal/priority"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: crawler <url>")
		os.Exit(1)
	}

	httpFetcher := fetcher.NewHTTPFetcher()
	// Priority + politeness scheduling lives in the frontier.
	priorityPoliteFrontier := frontier.NewPriorityPolite()
	memoryDeduper := dedupe.NewMemory()
	htmlParser := parser.NewHTMLParser()
	basicFilter := filter.NewBasicFilter()
	prioritizer := priority.NewSimple()

	for _, url := range os.Args[1:] {
		// Seed URLs default to low priority before assignment.
		req := model.CrawlRequest{URL: url, Priority: model.PriorityLow}
		req = prioritizer.Assign(req)
		priorityPoliteFrontier.Push(req)
	}

	crawlEngine := engine.New(
		httpFetcher,
		priorityPoliteFrontier,
		memoryDeduper,
		htmlParser,
		basicFilter,
		prioritizer,
	)

	if err := crawlEngine.Run(); err != nil {
		fmt.Println("crawl error:", err)
		os.Exit(1)
	}
}
