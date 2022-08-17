package cache_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"chen.com/distributecache/cache"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	a := assert.New(t)
	go cache.StartServer()
	time.Sleep(time.Second * 2)
	l, _ := net.Dial("tcp", "0.0.0.0:8086")
	conn := cache.NewConnection(l)
	go func() {

		for {
			conn.ReadFromLc()
			for res := range conn.DataChannel {
				rcmd, _ := cache.ParseCommand(res)
				if rcmd.Type == cache.GetOk {
					log.Printf("cmd is %v\n", rcmd)
					a.Equal(int(rcmd.Value.(float64)), 2)
				}
				if res == nil {
					break
				}
			}
		}
	}()
	for i := 0; i < 20; i++ {
		cmd1 := cache.NewCommand(cache.Add, fmt.Sprintf("test%d", i), i)
		b, _ := json.Marshal(cmd1)
		conn.WriteAdapter(b)
	}
	cmd := cache.NewCommand(cache.Update, "test1", 2)
	b, _ := json.Marshal(cmd)
	conn.WriteAdapter(b)
	cmd = cache.NewCommand(cache.Get, "test1", nil)
	b, _ = json.Marshal(cmd)
	conn.WriteAdapter(b)
	time.Sleep(time.Second * 2)
}
func TestTcp(t *testing.T) {
	l, _ := net.Listen("tcp", "0.0.0.0:8089")
	go func() {
		for {
			c, _ := l.Accept()
			for {
				b := make([]byte, 1024)
				len, _ := c.Read(b)
				log.Printf("str: %v\n", string(b[:len]))
			}
		}
	}()
	time.Sleep(time.Second * 1)
	c, _ := net.Dial("tcp", "localhost:8089")
	w := bufio.NewWriter(c)
	for i := 0; i < 20; i++ {
		w.Write([]byte(fmt.Sprintf("start%dend", i)))
		w.Flush()
		time.Sleep(time.Millisecond * 10)
	}
	time.Sleep(time.Second * 2)
}
