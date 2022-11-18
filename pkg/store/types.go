package store

type ISitemapStore interface {
	AddToSitemap(string, []string) error
	GetSitemap(string) (map[string]bool, error)
	SitemapExists(string) bool
	AddProgressToSitemap(string, string, []string)
	WaitAndGetSitemap(string) chan map[string]bool
}
