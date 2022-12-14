package crawler

import (
	"io"
	"log"
	"monzo-crawler/pkg/fetcher"
	"monzo-crawler/pkg/queue"
	"sync"

	"github.com/temoto/robotstxt"
)

var pLogger *log.Logger

type processor struct {
	depth   int
	writer  queue.IWriter
	fetcher fetcher.IFetcher
	fetched sync.Map
}

func NewProcessor(logOutput io.Writer, depth int, writer queue.IWriter, fetcher fetcher.IFetcher) IProcessor {
	pLogger = log.New(logOutput, "[processor]", log.LstdFlags)
	return &processor{
		depth:   depth,
		writer:  writer,
		fetcher: fetcher,
		fetched: sync.Map{},
	}
}

func (p *processor) Process(e queue.Element) {
	if e.GetType() != FETCH_URL {
		pLogger.Printf("Skip event. Only process FETCH_URL events. Event: %v\n", e.GetData())
		return
	}
	data, ok := e.GetData().(*queue.FetchElementData)
	if !ok {
		pLogger.Printf("Invalid data type. Expected *queue.ElementData. Got %v\n", data)
		return
	}
	toFetch := CreateToFetchUrl(data.CurUrl, data.Path)
	fetchedElementData := &queue.FetchedElementData{
		Urls:    []string{},
		Depth:   data.Depth,
		BaseUrl: data.BaseUrl,
		Path:    data.Path,
		CurUrl:  toFetch.String(),
		Robots:  data.Robots,
	}
	elementType := FETCHED_URL_ERROR

	defer func() {
		if elementType == FETCHED_URLS {
			p.fetched.Store(toFetch.String(), true)
		}
		if err := p.writer.Write(queue.NewFetchedQueueElement(fetchedElementData, e.GetBaseURL(), elementType)); err != nil {
			pLogger.Printf("Error while writing to fetched url stream. err: %v\n", err)
		}
	}()

	if !VerifySameDomain(data.BaseUrl, data.Path) {
		pLogger.Printf("Out of domain url")
		return
	}

	if fetched, ok := p.fetched.Load(toFetch.String()); ok && fetched.(bool) {
		pLogger.Printf("Already fetched. %v\n", toFetch)
		return
	}

	if data.Depth > p.depth {
		pLogger.Printf("Max depth reached. url %v\n", toFetch)
		return
	}
	robots, err := robotstxt.FromString(data.Robots)
	if err != nil {
		pLogger.Printf("Invalid robots.txt url %v\n", toFetch)
	}
	if robots != nil && !robots.TestAgent(toFetch.Path, "crawler") {
		pLogger.Printf("Cannot fetch url due to robots.txt restriction, url %v\n", toFetch)
		return
	}

	childUrls, err := p.fetcher.FetchChildURLs(toFetch.String())
	if err != nil {
		pLogger.Printf("Error while fetching url. Url: %v, err: %v\n", toFetch, err)
		return
	}
	elementType = FETCHED_URLS
	fetchedElementData.Urls = VerifyValidUrls(data.BaseUrl, childUrls)
}
