package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGetSitemap(t *testing.T) {
	s := NewSitemapStore()

	cases := []struct {
		name   string
		add    []string
		url    string
		expect map[string]bool
		err    bool
		before func()
	}{
		{
			name: "happy case",
			add:  []string{"1"},
			url:  "url",
			expect: map[string]bool{
				"1": true,
			},
			err:    false,
			before: func() {},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s.AddToSitemap(tc.url, tc.add)
			sm, e := s.GetSitemap(tc.url)
			if tc.err {
				assert.Error(t, e)
				return
			}
			assert.Equal(t, tc.expect, sm)
		})
	}
}

func TestSitemapInProcess(t *testing.T) {
	base := "base"
	cases := []struct {
		name   string
		add    []string
		url    string
		expect map[string]bool
		err    bool
		exists bool
		before func(ISitemapStore)
	}{
		{
			name: "in process",
			add:  []string{"1"},
			url:  "url",
			expect: map[string]bool{
				"1": true,
			},
			exists: false,
			err:    true,
			before: func(s ISitemapStore) {
				s.AddProgressToSitemap(base, "", []string{"/"})
			},
		},
		{
			name: "process finished",
			add:  []string{"1"},
			url:  "url",
			expect: map[string]bool{
				"1": true,
			},
			exists: true,
			err:    false,
			before: func(s ISitemapStore) {
				s.AddToSitemap(base, []string{})
				s.AddProgressToSitemap(base, "", []string{"1"})
				s.AddToSitemap(base, []string{"1"})
				s.AddProgressToSitemap(base, "1", []string{})
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSitemapStore()
			tc.before(s)
			assert.Equal(t, tc.exists, s.SitemapExists(base))
			if !tc.exists {
				return
			}
			res := <-s.WaitAndGetSitemap(base)
			assert.Equal(t, tc.expect, res)
		})
	}
}
