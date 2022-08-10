package data

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"chen.com/distributecache/config"
)

type DataSource struct {
	mux           sync.Mutex
	storeFilePath string
	f             *os.File
}

func NewDataSource() *DataSource {
	conf := config.GetConfig()
	ds := &DataSource{
		storeFilePath: fmt.Sprintf("%s/distribute_cache", conf.AppBaseDir),
	}
	var flag int
	if _, err := os.Stat(ds.storeFilePath); errors.Is(err, os.ErrNotExist) {
		flag = os.O_RDWR | os.O_CREATE
	} else {
		flag = os.O_RDWR
	}
	var err error
	ds.f, err = os.OpenFile(ds.storeFilePath, flag, 0644)
	if err != nil {
		log.Printf("error opening store file %s: %v", ds.storeFilePath, err)
	}
	return ds
}
func (ds *DataSource) ReadFromFile() (map[string]any, error) {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	f, err := os.OpenFile(ds.storeFilePath, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error opening store file")
		return nil, err
	}
	b := make([]byte, 1024)
	f.Read(b)
	memstr, _ := ds.decode(b)
	m := make(map[string]any)
	m[memstr] = ""
	return m, nil
}
func (ds *DataSource) WriteToFile(m map[string]any) error {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	for key, val := range m {
		b, _ := ds.encode(key)
		_, err := ds.f.Write(b)
		ds.addLine(ds.f)
		if err != nil {
			log.Printf("write file error: %v", err)
		}
		b, _ = ds.encode(val)
		_, err = ds.f.Write(b)
		ds.addLine(ds.f)
		if err != nil {
			log.Printf("write file error: %v", err)
		}
	}
	return nil
}
func (ds *DataSource) Close() {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	os.Remove(ds.storeFilePath)
}
func (ds *DataSource) encode(memstr any) ([]byte, error) {
	s, _ := json.Marshal(memstr)
	b := make([]byte, hex.EncodedLen(len(s)))
	hex.Encode(b, s)
	return b, nil
}
func (ds *DataSource) decode(b []byte) (string, error) {
	return string(b), nil
}
func (ds *DataSource) addLine(f *os.File) {
	str := "\r\n"
	f.Write([]byte(str))
}
