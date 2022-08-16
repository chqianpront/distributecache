package main

import (
	"log"

	"chen.com/distributecache/cache"
	"chen.com/distributecache/platform"
)

func main() {
	log.Printf("starting cache server")
	go platform.StartServ()
	cache.StartServer()
}
