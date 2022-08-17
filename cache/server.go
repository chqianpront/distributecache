package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"chen.com/distributecache/config"
)

const (
	magicNumber = 0x123456
)

var store CacheStore

func StartServer() {
	store = *NewStroe()
	conf := config.GetConfig()
	addr := fmt.Sprintf("%s:%d", conf.Addr, conf.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("error starting server: %v\n", err)
		os.Exit(-1)
	}
	log.Printf("listening on %s\n", addr)
	for {
		log.Printf("waiting for connection")
		conn, _ := l.Accept()
		connection := NewConnection(conn)
		go handleConn(connection)
	}
}
func handleConn(connection *Connection) {
	for {
		connection.SetReadDeadline(time.Now().Add(time.Second * 60))
		connection.ReadFromLc()
		for buf := range connection.DataChannel {
			if buf == nil {
				break
			}
			cmd, err := ParseCommand(buf)
			// log.Printf("cmd: %v\n", cmd)
			if err != nil {
				log.Printf("error parsing command: %v\n", err)
				continue
			}
			var bs []byte
			switch cmd.Type {
			case Add:
				store.Add(cmd.Key, cmd.Value)
				bs, _ = json.Marshal(NewCommand(Ok, "", nil))
			case Update:
				store.Put(cmd.Key, cmd.Value)
				bs, _ = json.Marshal(NewCommand(Ok, "", nil))
			case Delete:
				store.Delete(cmd.Key)
				bs, _ = json.Marshal(NewCommand(Ok, "", nil))
			case Get:
				gval, _ := store.Get(cmd.Key)
				bs, _ = json.Marshal(NewCommand(GetOk, cmd.Key, gval))
			case Ping:
				bs, _ = json.Marshal(NewCommand(Pong, "", nil))
			}
			connection.WriteAdapter(bs)
		}
	}
}
