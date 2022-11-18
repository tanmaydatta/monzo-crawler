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
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

var (
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
	if err := viper.UnmarshalKey(env, &config.Conf); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	config.Conf.Env = env
	LogFile := os.Stdout
	var err error
	if config.Conf.LogFile != "" {
		LogFile, err = os.OpenFile(config.Conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("error opening log file: %v", err))
		}
	}
	if config.Conf.LogFile != "" {
		config.Conf.LogFile = "log"
	}
	out := filepath.Join(config.Conf.OutFile, fmt.Sprintf("out_%v", time.Now().Unix()))
	if _, err := os.Stat("/path/to/your-file"); os.IsNotExist(err) {
		_ = os.MkdirAll(config.Conf.OutFile, 0700)
	}

	OutFile, err = os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening out file: %v", err))
	}

	fetcher := fetcher.NewFetcher(LogFile)
	store := store.NewSitemapStore()
	fetchQ := queue.NewQueue("crawler")
	fetchedQ := queue.NewQueue("crawler")
	Reader = queue.NewReader(fetchQ)
	Crawler = crawler.InitAndNewCrawler(LogFile, store, queue.NewReader(fetchedQ), queue.NewWriter(fetchQ))
	Processor = crawler.NewProcessor(LogFile, config.Conf.MaxDepth, queue.NewWriter(fetchedQ), fetcher)
}
