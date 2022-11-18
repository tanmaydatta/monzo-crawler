package fetcher

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestFetchURL(t *testing.T) {
	f := &fetcher{}

	r := mux.NewRouter()
	body := ""
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(body))
	})
	go http.ListenAndServe(":8181", r)

	cases := []struct {
		url      string
		expected []string
		err      bool
		name     string
		body     string
	}{
		{url: "s", expected: []string{}, err: true, name: "invalid URL for parser"},
		{url: "http://:8181", expected: []string{}, err: false, name: "non html body", body: "hello world"},
		{
			url:      "http://:8181",
			expected: []string{"foo", "/bar/baz"},
			err:      false,
			name:     "non html body",
			body:     `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body = tc.body
			urls, err := f.FetchURL(tc.url)
			assert.Equal(t, tc.expected, urls)
			assert.Equal(t, tc.err, err != nil)
		})
	}
}
