package fetcher

type IFetcher interface {
	FetchChildURLs(url string) ([]string, error)
	FetchRobotsTxt(u string) string
}
