package cache

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"chen.com/distributecache/config"
)

type MetaData struct {
	Pointer  int64
	DataSize int
	Line     int
	deleted  int
}
type DataSource struct {
	mux           sync.Mutex
	storeFilePath string
	f             *os.File
	pointer       int64
}

func NewDataSource() *DataSource {
	conf := config.GetConfig()
	ds := &DataSource{
		storeFilePath: fmt.Sprintf("%s/distribute_cache", conf.AppBaseDir),
	}
	var (
		flag        int
		pointerFlag int
	)

	if _, err := os.Stat(ds.storeFilePath); errors.Is(err, os.ErrNotExist) {
		flag = os.O_RDWR | os.O_CREATE
		pointerFlag = io.SeekStart
	} else {
		flag = os.O_RDWR
		pointerFlag = io.SeekEnd
	}
	var err error
	ds.f, err = os.OpenFile(ds.storeFilePath, flag, 0644)
	ds.pointer, _ = ds.f.Seek(0, pointerFlag)
	if err != nil {
		log.Printf("error opening store file %s: %v", ds.storeFilePath, err)
	}
	return ds
}
func (ds *DataSource) RestoreFromFile() (map[CacheKey]CacheValue, error) {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	m := make(map[CacheKey]CacheValue)
	sc := bufio.NewScanner(ds.f)
	num := 0
	meta := new(MetaData)
	var key CacheKey
	var kvstr string
	var keystr string
	var value any
	for sc.Scan() {
		line := sc.Bytes()
		switch num {
		case 0:
			str, _ := ds.decode(line)
			json.Unmarshal([]byte(str), &meta)
			num++
		case 1:
			kv, _ := ds.decode(line)
			kvstr = string(kv)
			kvArr := strings.Split(kvstr, "\r\n")
			keystr = kvArr[0]
			value = kvArr[1]
			num = 0
			metaStr, _ := json.Marshal(meta)
			v := CacheValue{
				meta:  string(metaStr),
				key:   CacheKey(keystr),
				value: value,
			}
			if meta.deleted == 0 {
				m[key] = v
			}
		}
	}
	ds.pointer, _ = ds.f.Seek(0, io.SeekEnd)
	return m, nil
}
func (ds *DataSource) WriteToFile(key string, value any, meta string, modifyType string) (MetaData, error) {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	metaData := new(MetaData)
	switch modifyType {
	case "create":
		metaData.Pointer = ds.pointer
		metaData.Line = 2
		metaData.deleted = 0
		wStr := fmt.Sprintf("%s\r\n%s\r\n", key, value)
		wb, _ := ds.encode([]byte(wStr))
		metaData.DataSize = len(wb)
		mb, _ := json.Marshal(metaData)
		mbd, _ := ds.encode(mb)
		ds.write(mbd, ds.pointer)
		ds.write([]byte("\r\n"), ds.pointer)
		ds.write(wb, ds.pointer)
		ds.write([]byte("\r\n"), ds.pointer)
		return *metaData, nil
	case "update":
		json.Unmarshal([]byte(meta), metaData)
		metaData.deleted = 1
		mb, _ := json.Marshal(metaData)
		mdb, _ := ds.encode(mb)
		ds.write(mdb, metaData.Pointer)

		wStr := fmt.Sprintf("%s\r\n%s\r\n", key, value)
		wb1, _ := ds.encode([]byte(wStr))
		metaData.DataSize = len(wb1)
		metaData.deleted = 0
		metaData.Pointer, _ = ds.f.Seek(0, os.SEEK_END)
		mb1, _ := json.Marshal(metaData)
		mbd1, _ := ds.encode(mb1)
		ds.write(mbd1, ds.pointer)
		ds.write(wb1, ds.pointer)
		ds.write([]byte("\r\n"), ds.pointer)
		return *metaData, nil

	case "delete":
		json.Unmarshal([]byte(meta), &metaData)
		metaData.deleted = 1
		mb, _ := json.Marshal(metaData)
		mdb, _ := ds.encode(mb)
		ds.f.WriteAt(mdb, metaData.Pointer)
		return MetaData{}, nil
	default:
		return MetaData{}, fmt.Errorf("unsupported modify type: %v", modifyType)
	}
}
func (ds *DataSource) Close() {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	os.Remove(ds.storeFilePath)
}
func (ds *DataSource) encode(tb []byte) ([]byte, error) {
	b := make([]byte, hex.EncodedLen(len(tb)))
	hex.Encode(b, tb)
	return b, nil
}
func (ds *DataSource) decode(b []byte) ([]byte, error) {
	rb := make([]byte, hex.DecodedLen(len(b)))
	hex.Decode(rb, b)
	return rb, nil
}
func (ds *DataSource) addLine(f *os.File) {
	pos, _ := ds.f.Seek(0, os.SEEK_END)
	str := "\r\n"
	ds.write([]byte(str), pos)
}
func (ds *DataSource) append(b []byte) error {
	pos, _ := ds.f.Seek(0, os.SEEK_END)
	err := ds.write(b, pos)
	if err != nil {
		log.Printf("error writing: err=%v", err)
		return err
	}
	ds.addLine(ds.f)
	return nil
}
func (ds *DataSource) write(b []byte, pos int64) error {
	wl, err := ds.f.WriteAt(b, pos)
	if err != nil {
		log.Printf("error writing: err=%v", err)
		return err
	}
	ds.pointer += int64(wl)
	return nil
}
