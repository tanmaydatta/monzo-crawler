package crawler

import (
	"monzo-crawler/pkg/queue"
)

type ICrawler interface {
	StartCrawl(string) error
	WaitAndGetSitemap(string) (map[string]bool, error)
}

type IProcessor interface {
	Process(queue.Element)
}

const (
	FETCH_URL         string = "fetch_url"
	FETCHED_URLS             = "fetched_urls"
	FETCHED_URL_ERROR        = "fetched_url_error"
)
