package fetcher

import (
	"io"
	"net/http"
	"web-crawler/internal/model"
)

type HTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{},
	}
}

func (f *HTTPFetcher) Fetch (req model.CrawlRequest) (*model.FetchResult, error) {

	resp, err := f.client.Get(req.URL)
	if err != nil {
		return  nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, err
	}

	return  &model.FetchResult{
		URL: req.URL,
		StatusCode: resp.StatusCode,
		Body: body,
	}, nil
}