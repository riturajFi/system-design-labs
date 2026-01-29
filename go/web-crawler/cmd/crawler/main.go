package main

import (
	"fmt"
	"os"
	"web-crawler/internal/dedupe"
	"web-crawler/internal/engine"
	"web-crawler/internal/fetcher"
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
	fifoFrontier := frontier.NewFIFO()
	memoryDeduper := dedupe.NewMemory()
	htmlParser := parser.NewHTMLParser()

	for _, url := range os.Args[1:] {
		fifoFrontier.Push(model.CrawlRequest{URL: url})
	}

	crawlEngine := engine.New(
		httpFetcher,
		fifoFrontier,
		memoryDeduper,
		htmlParser,
	)

	if err := crawlEngine.Run(); err != nil {
		fmt.Println("crawl error:", err)
		os.Exit(1)
	}
}
