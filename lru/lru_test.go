package lru_test

import (
	"fmt"
	"testing"

	"chen.com/distributecache/lru"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	a := assert.New(t)
	l := lru.Newlru()
	for i := 0; i < 100; i++ {
		if i%10 == 0 {
			l.Get("test1")
		}
		l.Add(fmt.Sprintf("test%d", i), i)
	}
	v, _ := l.Get("test1")
	a.Equal(v, 1)
}
