package main

import (
	"fmt"
	"log"
	it "monzo-crawler/init"
	"monzo-crawler/pkg/crawler"
	"monzo-crawler/pkg/queue"
)

var logger *log.Logger

func main() {
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
	baseUrl := "http://monzo.com/"
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
		it.OutFile.Write([]byte(fmt.Sprintf("%v\n", crawler.CreateToFetchUrl(baseUrl, url))))
	}
	it.OutFile.Close()
}
