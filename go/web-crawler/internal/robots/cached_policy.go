package robots

import (
	"net/http"
	"net/url"
	"sync"
	"web-crawler/internal/model"

	"github.com/temoto/robotstxt"
)

// CachedPolicy caches robots.txt per host and defaults to allow on errors.
type CachedPolicy struct {
	mu     sync.Mutex
	cache  map[string]*robotstxt.RobotsData
	client *http.Client
}

// NewCached constructs a cached robots.txt policy with a shared HTTP client.
func NewCached() *CachedPolicy {
	return &CachedPolicy{
		cache:  make(map[string]*robotstxt.RobotsData),
		client: &http.Client{},
	}
}

// Allowed checks robots.txt for the request host and path, caching per host.
func (p *CachedPolicy) Allowed(req model.CrawlRequest) bool {
	u, err := url.Parse(req.URL)
	if err != nil || u.Host == "" {
		return false
	}

	p.mu.Lock()
	data, ok := p.cache[u.Host]
	p.mu.Unlock()

	if !ok {
		data = p.fetchRobots(u)
		p.mu.Lock()
		p.cache[u.Host] = data
		p.mu.Unlock()
	}

	group := data.FindGroup("*")
	return group.Test(u.Path)
}

func (p *CachedPolicy) fetchRobots(u *url.URL) *robotstxt.RobotsData {
	resp, err := p.client.Get(u.Scheme + "://" + u.Host + "/robots.txt")
	if err != nil {
		data, _ := robotstxt.FromStatusAndBytes(200, nil)
		return data
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		fallback, _ := robotstxt.FromStatusAndBytes(resp.StatusCode, nil)
		return fallback
	}
	return data
}
