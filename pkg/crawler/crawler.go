package crawler

import (
	"io"
	"log"
	"monzo-crawler/pkg/queue"
	"monzo-crawler/pkg/store"
	urlpkg "net/url"
)

var cLogger *log.Logger

type crawler struct {
	siteMapStore store.ISitemapStore
	reader       queue.IReader
	writer       queue.IWriter
}

func (c *crawler) StartCrawl(url string) error {
	pUrl, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}
	base, _ := urlpkg.Parse("/")
	url = pUrl.ResolveReference(base).String()
	if c.siteMapStore.SitemapExists(url) || c.siteMapStore.SitemapInProcess(url) {
		return nil
	}
	c.siteMapStore.AddToSitemap(url, []string{})
	c.siteMapStore.AddProgressToSitemap(url, "", []string{"/"})
	return c.writer.Write(queue.NewFetchQueueElement(&queue.FetchElementData{
		Path: "/", BaseUrl: url, CurUrl: url, Depth: 1,
	}, url, FETCH_URL))
}

func (c *crawler) WaitAndGetSitemap(url string) (map[string]bool, error) {
	if c.siteMapStore.SitemapExists(url) {
		return c.siteMapStore.GetSitemap(url)
	}
	<-c.siteMapStore.WaitAndGetSitemap(url)
	// res, err := c.siteMapStore.GetSitemap(url)
	// if err != nil {
	// 	return nil, err
	// }
	return c.siteMapStore.GetSitemap(url)
	// ret := make(map[string]bool)
	// for k, v := range res {
	// 	CreateToFetchUrl()
	// }
}

func (c *crawler) processFetchedURLs() {
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
			for _, url := range data.Urls {
				if !VerifySameDomain(data.BaseUrl, url) {
					continue
				}
				urls = append(urls, url)
				go c.writer.Write(queue.NewFetchQueueElement(&queue.FetchElementData{
					Path:    url,
					BaseUrl: baseURL,
					CurUrl:  CreateToFetchUrl(data.CurUrl, url),
					Depth:   data.Depth + 1,
				}, baseURL, FETCH_URL))
			}

			c.siteMapStore.AddToSitemap(baseURL, urls)
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
	}
	go c.processFetchedURLs()
	return c
}
