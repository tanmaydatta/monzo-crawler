package main

import (
	"fmt"
	"log"
	it "monzo-crawler/init"
	"monzo-crawler/pkg/queue"
	"net/url"
	"os"
)

var logger *log.Logger

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		panic("No base url provided in args")
	}
	baseUrl, err := getBaseURL(args[0])
	if err != nil {
		panic(fmt.Errorf("Invalid argument %v", err))
	}
	it.Init()
	for i := 0; i < it.Conf.WorkersToFetchURL; i++ {
		go func() {
			for {
				ele, err := it.Reader.Read()
				if err == queue.EOF {
					return
				}
				if err != nil {
					logger.Println(err)
				}
				it.Processor.Process(ele)
			}
		}()
	}
	if err := it.Crawler.StartCrawl(baseUrl); err != nil {
		panic(fmt.Sprintf("couldn't start crawling %v", err))
	}
	urls, err := it.Crawler.WaitAndGetSitemap(baseUrl)
	if err != nil {
		fmt.Printf("Error occurred: %+v", err)
		return
	}
	it.OutFile.Write([]byte(fmt.Sprintf("Found %+v urls for %s\n========================================\n", len(urls), baseUrl)))
	for url := range urls {
		it.OutFile.Write([]byte(fmt.Sprintf("%v\n", url)))
	}
	it.OutFile.Close()
}

func getBaseURL(u string) (string, error) {
	in, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	home, _ := url.Parse("/")
	return in.ResolveReference(home).String(), nil
}
