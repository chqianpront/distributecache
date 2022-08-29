package raft

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Server struct {
	net.Conn
	C         chan []byte
	readEndCh chan bool
}

func NewServer(c net.Conn) *Server {
	p := &Server{
		Conn:      c,
		C:         make(chan []byte, 10),
		readEndCh: make(chan bool),
	}
	return p
}
func (p *Server) Write(b []byte) error {
	l := len(b)
	magicNum := make([]byte, 4)
	binary.BigEndian.PutUint32(magicNum, 0x123456)
	lenNum := make([]byte, 2)
	binary.BigEndian.PutUint16(lenNum, uint16(l))
	packetBuf := bytes.NewBuffer(magicNum)
	packetBuf.Write(lenNum)
	packetBuf.Write(b)
	_, err := p.Conn.Write(packetBuf.Bytes())
	return err
}
func (p *Server) Read() error {
	result := bytes.NewBuffer(nil)
	var buf [65542]byte // 由于 标识数据包长度 的只有两个字节 故数据包最大为 2^16+4(魔数)+2(长度标识)
	n, err := p.Conn.Read(buf[0:])
	result.Write(buf[0:n])
	if err != nil {
		return err
	} else {
		scanner := bufio.NewScanner(result)
		scanner.Split(p.splitFunc)
		for scanner.Scan() {
			fmt.Println("recv:", string(scanner.Bytes()[6:]))
			p.C <- scanner.Bytes()[6:]
		}
	}
	result.Reset()
	p.readEndCh <- true
	return nil
}
func (p *Server) Close() {
	close(p.readEndCh)
	close(p.C)
	p.Conn.Close()
}
func (p *Server) splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !atEOF && len(data) > 6 && binary.BigEndian.Uint32(data[:4]) == 0x123456 {
		var l int16
		binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &l)
		pl := int(l) + 6
		if pl <= len(data) {
			return pl, data[:pl], nil
		}
	}
	return
}
