package frontier

import "web-crawler/internal/model"

type Frontier interface {
	Push(req model.CrawlRequest)
	Pop() (model.CrawlRequest, bool)
	Len() int
	Done(req model.CrawlRequest)
}
