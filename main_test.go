package main_test

import (
	"log"
	"testing"
	"time"
)

func TestDistibuteCache(t *testing.T) {
	go func() {
		time.Sleep(time.Second * 10)
		log.Printf("exists")
	}()
	log.Printf("start")

}
