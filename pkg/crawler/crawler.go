package crawler

import (
	"io"
	"log"
	"monzo-crawler/pkg/fetcher"
	"monzo-crawler/pkg/queue"
	"monzo-crawler/pkg/store"
	urlpkg "net/url"
)

var cLogger *log.Logger

type crawler struct {
	siteMapStore store.ISitemapStore
	reader       queue.IReader
	writer       queue.IWriter
	fetcher      fetcher.IFetcher
}

func (c *crawler) StartCrawl(url string) error {
	pUrl, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}
	base, _ := urlpkg.Parse("/")
	url = pUrl.ResolveReference(base).String()
	if c.siteMapStore.SitemapExists(url) {
		return nil
	}
	c.siteMapStore.AddToSitemap(url, []string{})
	c.siteMapStore.AddProgressToSitemap(url, "", []string{"/"})
	return c.writer.Write(queue.NewFetchQueueElement(&queue.FetchElementData{
		Path: "/", BaseUrl: url, CurUrl: url, Depth: 1, Robots: c.fetcher.FetchRobotsTxt(url),
	}, url, FETCH_URL))
}

func (c *crawler) WaitAndGetSitemap(url string) (map[string]bool, error) {
	if c.siteMapStore.SitemapExists(url) {
		return c.siteMapStore.GetSitemap(url)
	}
	<-c.siteMapStore.WaitAndGetSitemap(url)
	return c.siteMapStore.GetSitemap(url)
}

func (c *crawler) processFetchedURLs() {
	defer func() {
		if err := recover(); err != nil {
			cLogger.Println("panic occurred:", err)
		}
	}()
	var e queue.Element
	var err error
	for e, err = c.reader.Read(); err != queue.EOF; e, err = c.reader.Read() {
		if err != nil {
			cLogger.Println(err)
			continue
		}
		baseURL := e.GetBaseURL()
		switch e.GetType() {
		case FETCH_URL:
		case FETCHED_URL_ERROR:
			data := e.GetData().(*queue.FetchedElementData)
			c.siteMapStore.AddProgressToSitemap(baseURL, data.Path, []string{})
		case FETCHED_URLS:
			data := e.GetData().(*queue.FetchedElementData)
			urls := make([]string, 0, len(data.Urls))
			urlsToFetch := make([]string, 0, len(data.Urls))
			for _, url := range data.Urls {
				if !VerifySameDomain(data.BaseUrl, url) {
					continue
				}
				toFetch := CreateToFetchUrl(data.CurUrl, url)
				urlsToFetch = append(urlsToFetch, toFetch.String())
				urls = append(urls, url)

				go c.writer.Write(queue.NewFetchQueueElement(&queue.FetchElementData{
					Path:    url,
					BaseUrl: baseURL,
					CurUrl:  toFetch.String(),
					Depth:   data.Depth + 1,
					Robots:  data.Robots,
				}, baseURL, FETCH_URL))
			}

			c.siteMapStore.AddToSitemap(baseURL, urlsToFetch)
			c.siteMapStore.AddProgressToSitemap(baseURL, data.Path, urls)
		}
	}
}

func InitAndNewCrawler(logOutput io.Writer, sitemapStore store.ISitemapStore, reader queue.IReader, writer queue.IWriter) ICrawler {
	cLogger = log.New(logOutput, "[crawler]", log.LstdFlags)

	c := &crawler{
		reader:       reader,
		writer:       writer,
		siteMapStore: sitemapStore,
		fetcher:      fetcher.NewFetcher(logOutput),
	}
	go c.processFetchedURLs()
	return c
}
