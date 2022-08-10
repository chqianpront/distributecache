package cache_test

import (
	"fmt"
	"testing"

	"chen.com/distributecache/cache"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	a := assert.New(t)
	s := cache.NewStroe()
	for i := 0; i < 1000; i++ {
		s.Add(fmt.Sprintf("test%d", i), i)
	}
	val, _ := s.Get("test2")
	a.Equal(val, 2)
	s.Destroy()
}
func TestDestroy(t *testing.T) {
	s := cache.NewStroe()
	s.Add("test", "test")
	s.Destroy()
}
