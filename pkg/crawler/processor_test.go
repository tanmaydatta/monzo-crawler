package crawler

import (
	"log"
	fetcherM "monzo-crawler/pkg/fetcher/mocks"
	"monzo-crawler/pkg/queue"
	queueM "monzo-crawler/pkg/queue/mocks"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessElement(t *testing.T) {
	mockWriter := queueM.IWriter{}
	mockFetcher := fetcherM.IFetcher{}
	ele := queueM.Element{}
	defer func() {
		mockFetcher.AssertExpectations(t)
		mockWriter.AssertExpectations(t)
		ele.AssertExpectations(t)
	}()
	pLogger = log.New(os.Stderr, "[processor]", log.LstdFlags)
	p := &processor{
		depth:   1,
		writer:  &mockWriter,
		fetcher: &mockFetcher,
		fetched: sync.Map{},
	}
	baseUrl := "http://monzo.com/"
	cases := []struct {
		name   string
		before func()
		ele    queue.Element
	}{
		{
			name: "Skip element",
			before: func() {
				p.fetched = sync.Map{}
			},
			ele: queue.NewFetchedQueueElement(&queue.FetchedElementData{}, baseUrl, FETCHED_URLS),
		},
		{
			name: "Invalid element data",
			before: func() {
				p.fetched = sync.Map{}
				ele.On("GetData").Return("").Once()
				ele.On("GetType").Return(FETCH_URL).Once()
			},
			ele: &ele,
		},
		{
			name: "Out of domain",
			before: func() {
				p.fetched = sync.Map{}
				mockWriter.On("Write", mock.Anything).Return(nil).Once()
			},
			ele: queue.NewFetchQueueElement(&queue.FetchElementData{
				BaseUrl: baseUrl,
				CurUrl:  baseUrl,
				Path:    "http://google.com/",
				Depth:   1,
			}, baseUrl, FETCH_URL),
		},
		{
			name: "Already fetched",
			before: func() {
				toFetch := CreateToFetchUrl(baseUrl, "/")
				p.fetched = sync.Map{}
				p.fetched.Store(toFetch.String(), true)
				mockWriter.On("Write", mock.Anything).Return(nil).Once()
			},
			ele: queue.NewFetchQueueElement(&queue.FetchElementData{
				BaseUrl: baseUrl,
				CurUrl:  baseUrl,
				Path:    "/",
				Depth:   1,
			}, baseUrl, FETCH_URL),
		},
		{
			name: "Exceeded max depth",
			before: func() {
				p.fetched = sync.Map{}
				mockWriter.On("Write", mock.Anything).Return(nil).Once()
			},
			ele: queue.NewFetchQueueElement(&queue.FetchElementData{
				BaseUrl: baseUrl,
				CurUrl:  baseUrl,
				Path:    "/",
				Depth:   2,
			}, baseUrl, FETCH_URL),
		},
		{
			name: "disallow due to robots.txt",
			before: func() {
				p.fetched = sync.Map{}
				mockWriter.On("Write", mock.Anything).Return(nil).Once()
			},
			ele: queue.NewFetchQueueElement(&queue.FetchElementData{
				BaseUrl: baseUrl,
				CurUrl:  baseUrl,
				Path:    "/pay/",
				Depth:   1,
				Robots: `
				# robotstxt.org/

				User-agent: *
				Disallow: /docs/
				Disallow: /referral/
				Disallow: /-staging-referral/
				Disallow: /install/
				Disallow: /blog/authors/
				Disallow: /pay/
				`,
			}, baseUrl, FETCH_URL),
		},
		{
			name: "Call fetcher get error",
			before: func() {
				p.fetched = sync.Map{}
				toFetch := CreateToFetchUrl(baseUrl, "/")
				mockFetcher.On("FetchChildURLs", toFetch.String()).Return([]string{"child_url"}, assert.AnError).Once()
				mockWriter.On("Write", mock.Anything).Return(nil).Once()
			},
			ele: queue.NewFetchQueueElement(&queue.FetchElementData{
				BaseUrl: baseUrl,
				CurUrl:  baseUrl,
				Path:    "/",
				Depth:   1,
			}, baseUrl, FETCH_URL),
		},
		{
			name: "Call fetcher happy case",
			before: func() {
				p.fetched = sync.Map{}
				toFetch := CreateToFetchUrl(baseUrl, "/")
				mockFetcher.On("FetchChildURLs", toFetch.String()).Return([]string{"child_url"}, nil).Once()
				mockWriter.On("Write", mock.Anything).Run(func(args mock.Arguments) {
					e := args.Get(0).(queue.Element)
					assert.Equal(t, FETCHED_URLS, e.GetType())
				}).Return(nil).Once()
			},
			ele: queue.NewFetchQueueElement(&queue.FetchElementData{
				BaseUrl: baseUrl,
				CurUrl:  baseUrl,
				Path:    "/",
				Depth:   1,
			}, baseUrl, FETCH_URL),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			p.Process(tc.ele)
		})
	}
}
