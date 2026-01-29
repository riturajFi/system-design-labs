package main

import (
	"fmt"
	"os"
	"web-crawler/internal/engine"
	"web-crawler/internal/fetcher"
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
	crawlEnginer := engine.New(httpFetcher)

	if err := crawlEnginer.Run(seedUrl); err != nil {
		fmt.Println("crawl error:", err)
		os.Exit(1)
	}
}