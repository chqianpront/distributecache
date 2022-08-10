package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sync"

	"chen.com/distributecache/data"
	"chen.com/distributecache/lru"
)

type CacheStore struct {
	mux         sync.Mutex
	l           *lru.Lru
	ds          *data.DataSource
	m           map[string]any
	flushCh     chan int
	flushCircle int
}

func NewStroe() *CacheStore {
	cs := &CacheStore{
		l:           lru.Newlru(),
		ds:          data.NewDataSource(),
		m:           make(map[string]any),
		flushCh:     make(chan int),
		flushCircle: 0,
	}
	go func() {
		for v := range cs.flushCh {
			log.Printf("received channel %v", v)
			cs.mux.Lock()
			tm := make(map[string]any)
			cs.snapshot(&tm)
			cs.mux.Unlock()
			cs.ds.WriteToFile(tm)
		}
	}()
	return cs
}
func (cs *CacheStore) Add(key string, value any) error {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	if cs.l == nil {
		cs.l = lru.Newlru()
	}
	err := cs.l.Add(key, value)
	cs.m[key] = value
	if err != nil {
		log.Printf("add to store failed")
		return err
	}
	cs.flushCircle++
	if cs.flushCircle%100 == 0 {
		cs.flushCh <- 1
	}
	return nil
}
func (cs *CacheStore) Get(key string) (v any, err error) {
	var ret any
	if cs.l != nil {
		ret, err = cs.l.Get(key)
		if err == nil {
			return ret, nil
		}
		log.Printf("get from cache error: %v", err)
	}
	if ret, ok := cs.m[key]; ok {
		return ret, nil
	} else {
		return nil, fmt.Errorf("%s not found", key)
	}
}
func (cs *CacheStore) Put(key string, value any) error {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	if cs.l == nil {
		cs.l = lru.Newlru()
	}
	cs.l.Remove(key)
	cs.l.Add(key, value)
	if _, ok := cs.m[key]; ok {
		return fmt.Errorf("can not find cache key %s", key)
	}
	cs.m[key] = value
	return nil
}
func (cs *CacheStore) Delete(key string) {
	cs.mux.Lock()
	defer cs.mux.Unlock()
	if cs.l != nil {
		cs.l.Remove(key)
	}
	delete(cs.m, key)
}
func (cs *CacheStore) Destroy() {
	cs.ds.Close()
	close(cs.flushCh)

}
func (cs *CacheStore) snapshot(out any) {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(cs.m)
	gob.NewDecoder(buf).Decode(out)
}
