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
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: crawler <url>")
		os.Exit(1)
	}

	httpFetcher := fetcher.NewHTTPFetcher()
	// Polite frontier enforces one in-flight request per host.
	politeFrontier := frontier.NewPolite()
	memoryDeduper := dedupe.NewMemory()
	htmlParser := parser.NewHTMLParser()
	basicFilter := filter.NewBasicFilter()

	for _, url := range os.Args[1:] {
		politeFrontier.Push(model.CrawlRequest{URL: url})
	}

	crawlEngine := engine.New(
		httpFetcher,
		politeFrontier,
		memoryDeduper,
		htmlParser,
		basicFilter,
	)

	if err := crawlEngine.Run(); err != nil {
		fmt.Println("crawl error:", err)
		os.Exit(1)
	}
}
