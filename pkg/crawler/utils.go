package crawler

import (
	"net/url"
)

func VerifySameDomain(baseUrl, path string) bool {
	b, _ := url.Parse(baseUrl)
	u, _ := url.Parse(path)
	u.Fragment = ""
	b.Fragment = ""
	return !(u.Host != "" && u.Host != b.Host)
}

func CreateToFetchUrl(baseUrl, path string) *url.URL {

	b, _ := url.Parse(baseUrl)
	u, _ := url.Parse(path)
	u.Fragment = ""
	b.Fragment = ""
	// fmt.Printf("%v %v %v\n", b.ResolveReference(u).String(), baseUrl, path)
	return b.ResolveReference(u)
}

func VerifyValidUrls(baseUrl string, urls []string) []string {
	ret := make([]string, 0, len(urls))
	b, _ := url.Parse(baseUrl)
	for _, u := range urls {
		u2, err := url.Parse(u)
		if err != nil {
			continue
		}
		if u2, err = url.ParseRequestURI(b.ResolveReference(u2).String()); err != nil {
			continue
		}
		if u2.Scheme != "http" && u2.Scheme != "https" {
			continue
		}
		ret = append(ret, u)
	}
	return ret
}
