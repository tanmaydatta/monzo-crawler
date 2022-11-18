package fetcher

import (
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type fetcher struct {
}

type parser struct {
	childURLs []string
}

func (f *fetcher) FetchURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	return f.parseAndGetChildURLs(string(body))
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

func NewFetcher() IFetcher {
	return &fetcher{}
}
