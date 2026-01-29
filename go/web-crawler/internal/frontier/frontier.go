package frontier

import "web-crawler/internal/model"

type Frontier interface {
	Push(req model.CrawlRequest)
	Pop() (model.CrawlRequest, bool)
	Len() int
	// Done notifies the frontier that a request has finished processing.
	Done(req model.CrawlRequest)
}
