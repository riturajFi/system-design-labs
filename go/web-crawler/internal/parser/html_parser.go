package parser

import (
	"bytes"
	"net/url"
	"web-crawler/internal/model"

	"golang.org/x/net/html"
)

type HTMLParser struct{}

func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

func (p *HTMLParser) Parse(baseURL string, body []byte) ([]model.CrawlRequest, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	var results []model.CrawlRequest

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}
					abs := base.ResolveReference(u)
					// Default to low priority; prioritizer can override later.
					results = append(results, model.CrawlRequest{
						URL:      abs.String(),
						Priority: model.PriorityLow,
					})
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	return results, nil
}
