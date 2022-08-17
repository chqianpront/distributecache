package cache

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

type Connection struct {
	Id              string      `json:"id"`
	IsAuthenticated bool        `json:"is_authenticated"`
	DataChannel     chan []byte `json:"data_channel"`
	net.Conn
}

func NewConnection(c net.Conn) *Connection {
	return &Connection{
		IsAuthenticated: false,
		Conn:            c,
		DataChannel:     make(chan []byte),
	}
}
func (conn *Connection) ReadFromLc() {
	go func() {
		var res []byte
		result := bytes.NewBuffer(nil)
		var buf [65542]byte // 由于 标识数据包长度 的只有两个字节 故数据包最大为 2^16+4(魔数)+2(长度标识)
		n, err := conn.Conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
			} else {

			}
		} else {
			scanner := bufio.NewScanner(result)
			scanner.Split(packetSlitFunc)
			for scanner.Scan() {
				res = scanner.Bytes()[6:]
				conn.DataChannel <- res
			}
		}
		result.Reset()
		conn.DataChannel <- nil
	}()
}
func (conn *Connection) writeToLc(b []byte) error {
	_, err := conn.Conn.Write(b)
	return err
}
func (conn *Connection) WriteAdapter(b []byte) error {
	l := len(b)
	magicNum := make([]byte, 4)
	binary.BigEndian.PutUint32(magicNum, magicNumber)
	lenNum := make([]byte, 2)
	binary.BigEndian.PutUint16(lenNum, uint16(l))
	packetBuf := bytes.NewBuffer(magicNum)
	packetBuf.Write(lenNum)
	packetBuf.Write(b)

	return conn.writeToLc(packetBuf.Bytes())
}

func packetSlitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !atEOF && len(data) > 6 && binary.BigEndian.Uint32(data[:4]) == magicNumber {
		var l int16
		binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &l)
		pl := int(l) + 6
		if pl <= len(data) {
			return pl, data[:pl], nil
		}
	}
	return
}
