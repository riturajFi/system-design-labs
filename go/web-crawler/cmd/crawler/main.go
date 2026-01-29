package main

import (
	"fmt"
	"os"
	"web-crawler/internal/engine"
	"web-crawler/internal/fetcher"
	"web-crawler/internal/frontier"
	"web-crawler/internal/model"
)

func main() {

	if len(os.Args) < 2 {

		if len(os.Args) < 2 {
			fmt.Println("usage: crawler <url>")
			os.Exit(1)
		}
	}

	seedUrl := os.Args[1]

	httpFetcher := fetcher.NewHTTPFetcher()
	fifoFrontier := frontier.NewFIFO()

	for _, url := range os.Args[1:] {
		fifoFrontier.Push(model.CrawlRequest{URL: url})
	}

	crawlEnginer := engine.New(httpFetcher)

	if err := crawlEnginer.Run(seedUrl); err != nil {
		fmt.Println("crawl error:", err)
		os.Exit(1)
	}
}