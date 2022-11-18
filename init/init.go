package init

import (
	"fmt"
	"io"
	"monzo-crawler/config"
	"monzo-crawler/pkg/crawler"
	"monzo-crawler/pkg/fetcher"
	"monzo-crawler/pkg/queue"
	"monzo-crawler/pkg/store"
	"os"

	"github.com/spf13/viper"
)

var (
	Conf      config.Config
	Crawler   crawler.ICrawler
	Processor crawler.IProcessor
	Reader    queue.IReader
	LogFile   io.WriteCloser
	OutFile   io.WriteCloser
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	env := viper.GetString("env")
	if err := viper.UnmarshalKey(env, &Conf); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	Conf.Env = env
	LogFile := os.Stdout
	var err error
	if Conf.LogFile != "" {
		LogFile, err = os.OpenFile(Conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("error opening log file: %v", err))
		}
	}
	if Conf.LogFile != "" {
		Conf.LogFile = "out/out"
	}
	OutFile, err = os.OpenFile(Conf.OutFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening out file: %v", err))
	}

	fetcher := fetcher.NewFetcher(LogFile)
	store := store.NewSitemapStore()
	fetchQ := queue.NewQueue("crawler")
	fetchedQ := queue.NewQueue("crawler")
	Reader = queue.NewReader(fetchQ)
	Crawler = crawler.InitAndNewCrawler(LogFile, store, queue.NewReader(fetchedQ), queue.NewWriter(fetchQ))
	Processor = crawler.NewProcessor(LogFile, Conf.MaxDepth, queue.NewWriter(fetchedQ), fetcher)
}
