package cache_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"chen.com/distributecache/cache"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	a := assert.New(t)
	s := cache.NewStroe()
	for i := 0; i < 103; i++ {
		s.Add(fmt.Sprintf("test%d", i), i)
	}
	val, _ := s.Get("test1")
	a.Equal(val, 1)
	// time.Sleep(time.Second * 3)
	s.Destroy()
}
func TestDestroy(t *testing.T) {
	s := cache.NewStroe()
	s.Put("test3", 345)
	s.Destroy()
}
func TestWriteFile(t *testing.T) {
	f, _ := os.OpenFile("test", os.O_WRONLY|os.O_CREATE, 0644)
	pos, _ := f.Seek(0, os.SEEK_SET)
	f.WriteAt([]byte("test\r\n"), pos)
	pos, _ = f.Seek(0, os.SEEK_END)
	wl, err := f.WriteAt([]byte("heelo\r\n"), pos)
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Printf("write length is %d", wl)
}
