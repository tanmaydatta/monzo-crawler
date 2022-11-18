package store

import (
	"sync"
)

type sitemapStore struct {
	values    map[string]map[string]bool
	inProcess map[string]map[string]bool
	mu        sync.Mutex
	results   map[string]chan map[string]bool
}

func NewSitemapStore() ISitemapStore {
	return &sitemapStore{
		values:    make(map[string]map[string]bool),
		inProcess: make(map[string]map[string]bool),
		results:   make(map[string]chan map[string]bool),
		mu:        sync.Mutex{},
	}
}

func (s *sitemapStore) AddToSitemap(url string, urls []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.values[url]
	if !ok {
		val = map[string]bool{}
	}
	for _, url := range urls {
		val[url] = true
	}
	s.values[url] = val
	return nil
}

func (s *sitemapStore) GetSitemap(url string) (map[string]bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, _ := s.values[url]
	return val, nil
}

func (s *sitemapStore) SitemapExists(url string) bool {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	val, ok := s.inProcess[url]
	if !ok {
		return false
	}
	return len(val) == 0
}

func (s *sitemapStore) AddProgressToSitemap(baseUrl string, remove string, add []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var ok bool
	var val map[string]bool
	val, ok = s.inProcess[baseUrl]
	if !ok {
		val = map[string]bool{}
	}
	if _, ok = val[remove]; ok {
		delete(val, remove)
	}
	for _, url := range add {
		val[url] = true
	}
	s.inProcess[baseUrl] = val
	go func() {
		if s.SitemapExists(baseUrl) {
			sm, _ := s.GetSitemap(baseUrl)
			s.mu.Lock()
			if _, ok := s.results[baseUrl]; !ok {
				s.results[baseUrl] = make(chan map[string]bool)
			}
			res := s.results[baseUrl]
			s.mu.Unlock()
			res <- sm
		}
	}()
}

func (s *sitemapStore) WaitAndGetSitemap(url string) chan map[string]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.results[url]; !ok {
		s.results[url] = make(chan map[string]bool)
	}
	return s.results[url]
}
