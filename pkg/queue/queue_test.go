package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQ(t *testing.T) {
	q := NewQueue("q")
	baseURL := "base"
	cases := []struct {
		expect string
		err    bool
		name   string
		before func()
	}{
		{
			name:   "happy case",
			expect: baseURL,
			err:    false,
			before: func() {
				go q.EnQueue(NewFetchQueueElement(nil, baseURL, ""))
			},
		},
		{
			name:   "dequeue without enqueue",
			expect: baseURL,
			err:    true,
			before: func() {},
		},
		{
			name:   "dequeue after close",
			expect: baseURL,
			err:    true,
			before: func() {
				q.Close()
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			res := make(chan Element)
			go func() {
				e, _ := q.DeQueue()
				res <- e
			}()
			var e Element
			select {
			case e = <-res:
			case <-time.After(1 * time.Second):
			}
			if tc.err {
				assert.Nil(t, e)
				return
			}
			assert.NotNil(t, e)
			assert.Equal(t, tc.expect, e.GetBaseURL())
		})
	}
}
