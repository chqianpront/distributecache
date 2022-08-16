package cache_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"chen.com/distributecache/cache"
)

func TestServer(t *testing.T) {
	// a := assert.New(t)
	go cache.StartServer()
	time.Sleep(time.Second * 2)
	l, _ := net.Dial("tcp", "0.0.0.0:8086")
	conn := cache.NewConnection(l)
	go func() {

		for {
			// err := conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			// if err != nil {
			// 	break
			// }
			res := conn.ReadFromLc()
			log.Printf("res str: %v", string(res))
		}
	}()
	for i := 0; i < 20; i++ {
		cmd1 := cache.NewCommand(cache.Add, fmt.Sprintf("test%d", i), i)
		b, _ := json.Marshal(cmd1)
		conn.WriteAdapter(b)
	}
	cmd := cache.NewCommand(cache.Get, "test1", nil)
	b, _ := json.Marshal(cmd)
	conn.WriteAdapter(b)
	time.Sleep(time.Second * 2)
}
