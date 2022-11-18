package crawler

import "net/url"

func VerifySameDomain(baseUrl, path string) bool {
	b, _ := url.Parse(baseUrl)
	u, _ := url.Parse(path)
	u.Fragment = ""
	b.Fragment = ""
	return !(u.Host != "" && u.Host != b.Host)
}

func CreateToFetchUrl(baseUrl, path string) string {
	b, _ := url.Parse(baseUrl)
	u, _ := url.Parse(path)
	u.Fragment = ""
	b.Fragment = ""
	return b.ResolveReference(u).String()
}
