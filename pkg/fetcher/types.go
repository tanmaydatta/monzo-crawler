package fetcher

type IFetcher interface {
	FetchURL(url string) ([]string, error)
}
