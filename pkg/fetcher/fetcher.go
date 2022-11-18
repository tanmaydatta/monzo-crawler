package fetcher

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var fLogger *log.Logger

type fetcher struct {
	client *http.Client
}

type parser struct {
	childURLs []string
}

func (f *fetcher) FetchChildURLs(url string) ([]string, error) {
	defer func() {
		if err := recover(); err != nil {
			fLogger.Println("panic occurred:", err)
		}
	}()

	body, err := f.fetchURL(url)
	if err != nil {
		return []string{}, err
	}

	return f.parseAndGetChildURLs(body)
}

func (f *fetcher) fetchURL(url string) (string, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (f *fetcher) FetchRobotsTxt(u string) string {
	r, _ := url.Parse("robots.txt")
	base, _ := url.Parse(u)
	body, err := f.fetchURL(base.ResolveReference(r).String())
	if err != nil {
		fLogger.Println("error in fetching robots.txt", err)
	}
	return body
}

func (f *fetcher) parseAndGetChildURLs(body string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return []string{}, err
	}
	p := &parser{
		childURLs: []string{},
	}
	p.parse(doc)
	return p.childURLs, nil
}

func (p *parser) parse(node *html.Node) {
	if node == nil {
		return
	}
	for _, a := range node.Attr {
		if a.Key == "href" && node.Data == "a" {
			p.childURLs = append(p.childURLs, a.Val)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		p.parse(c)
	}
}

func NewFetcher(logOutput io.Writer) IFetcher {
	fLogger = log.New(logOutput, "[fetcher]", log.LstdFlags)
	client := &http.Client{Timeout: 2 * time.Second}
	return &fetcher{client: client}
}
