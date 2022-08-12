package cache_test

import (
	"log"
	"testing"

	"chen.com/distributecache/cache"
)

func TestCodec(t *testing.T) {
	d := cache.NewDataSource()
	d.WriteToFile("test", 345, "", "create")
	m, _ := d.RestoreFromFile()
	log.Printf("map: %v", m)
}
