package lru

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"chen.com/distributecache/config"
)

type Lru struct {
	root *list.Element
	tail *list.Element
	l    *list.List
	m    map[string]*list.Element
	size int
}
type Node struct {
	key   any
	value any
}

func Newlru() *Lru {
	conf := config.GetConfig()
	lru := &Lru{
		root: &list.Element{},
		tail: nil,
		l:    list.New(),
		m:    make(map[string]*list.Element),
		size: conf.LruSize,
	}
	return lru
}
func (lru *Lru) Len() int { return lru.l.Len() }
func (lru *Lru) Add(key, v any) error {
	var elm *list.Element
	valKey, err := lru.keyHash(key)
	if err != nil {
		log.Printf("get key hash error: %v", err)
		return err
	}
	if _, dup := lru.m[valKey]; dup {
		log.Printf("key already exists: %v", valKey)
		return errors.New("key already exists")
	}
	node := Node{key: key, value: v}
	if elm = lru.l.PushFront(node); elm == nil {
		return errors.New("lru not initalized")
	}
	lru.m[valKey] = elm
	lru.root = elm
	if lru.tail == nil {
		lru.tail = elm
	}
	if lru.Len() > lru.size {
		lru.removeOverSize()
	}
	return nil
}
func (lru *Lru) Get(key any) (v any, err error) {
	kHash, err := lru.keyHash(key)
	if err != nil {
		log.Printf("get key hash error: %v", err)
		return nil, err
	}
	if elm, ok := lru.m[kHash]; ok {
		return elm.Value.(Node).value, nil
	}
	return nil, fmt.Errorf("%v not found", key)
}
func (lru *Lru) Remove(key any) {
	delete(lru.m, key.(string))
}
func (lru *Lru) keyHash(key any) (string, error) {
	keystr := key.(string)
	d := []byte(keystr)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil)), nil
}
func (lru *Lru) removeOverSize() {
	for lru.Len() > lru.size {
		tmp := lru.tail
		lru.tail = lru.tail.Prev()
		lru.l.Remove(tmp)
	}
}
