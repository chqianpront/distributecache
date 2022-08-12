package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"chen.com/distributecache/lru"
)

type CacheKey string
type CacheValue struct {
	meta  string
	key   CacheKey
	value any
}
type CacheStore struct {
	mux sync.Mutex
	l   *lru.Lru
	ds  *DataSource
	m   map[CacheKey]CacheValue
}

func NewStroe() *CacheStore {
	cs := &CacheStore{
		l:  lru.Newlru(),
		ds: NewDataSource(),
		m:  make(map[CacheKey]CacheValue),
	}
	var err error
	cs.m, err = cs.ds.RestoreFromFile()
	if err != nil {
		log.Printf("Failed to new cache store: %v", err)
	}
	return cs
}
func (cs *CacheStore) Add(key CacheKey, value any) error {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	if cs.l == nil {
		cs.l = lru.Newlru()
	}
	err := cs.l.Add(string(key), value)
	if err != nil {
		log.Printf("add to store failed")
		return err
	}
	meta, err := cs.ds.WriteToFile(string(key), value, "", "create")
	if err != nil {
		log.Printf("add to store failed")
		return err
	}
	metaStr, _ := json.Marshal(meta)
	v := CacheValue{
		meta:  string(metaStr),
		key:   key,
		value: value,
	}
	cs.m[key] = v
	return nil
}
func (cs *CacheStore) Get(key CacheKey) (v any, err error) {
	var ret any
	if cs.l != nil {
		ret, err = cs.l.Get(string(key))
		if err == nil {
			return ret, nil
		}
		log.Printf("get from cache error: %v", err)
	}
	if ret, ok := cs.m[key]; ok {
		return ret.Value(), nil
	} else {
		return nil, fmt.Errorf("%s not found", key)
	}
}
func (cs *CacheStore) Put(key CacheKey, value any) error {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	if cs.l == nil {
		cs.l = lru.Newlru()
	}
	cs.l.Remove(string(key))
	cs.l.Add(string(key), value)
	v := cs.m[key]
	metaStr, _ := json.Marshal(v.MetaData())
	meta, err := cs.ds.WriteToFile(string(key), value, string(metaStr), "update")
	if err != nil {
		log.Printf("error writing update meta data: %v", err)
		return err
	}
	metaStr1, _ := json.Marshal(meta)
	newV := CacheValue{
		meta:  string(metaStr1),
		key:   key,
		value: value,
	}
	cs.m[key] = newV
	return nil
}
func (cs *CacheStore) Delete(key CacheKey) {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	if cs.l != nil {
		cs.l.Remove(string(key))
	}
	delete(cs.m, key)
}
func (cs *CacheStore) Destroy() {
	cs.ds.Close()
}
func (cv *CacheValue) Value() any {
	return cv.value
}
func (cv *CacheValue) Key() string {
	return string(cv.key)
}
func (cv *CacheValue) MetaData() *MetaData {
	var metaData MetaData
	json.Unmarshal([]byte(cv.meta), &metaData)
	return &metaData
}
