package gocache

import (
	"container/list"
)

 
type CacheLRU struct {
	maxBytes int64
	nbytes int64
	ll	*list.List //lru算法的调度队列
	cache map[string]*list.Element//存储数据的结构
	onEvicted func(key string,value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewCache(maxBytes int64,onEvicted func(key string,value Value)) *CacheLRU {
	return &CacheLRU{
		maxBytes: maxBytes,
		ll:list.New(),
		cache: make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *CacheLRU)Get(key string) (value Value,ok bool){
	if ele ,err := c.cache[key]; err{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value,true
	}
	return
}

func (c *CacheLRU)DeleteOldest(){
	ele := c.ll.Back()
	if ele != nil{
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache,kv.key)
		c.nbytes -= int64(len(kv.key))+ int64(kv.value.Len())
		if c.onEvicted != nil{
			c.onEvicted(kv.key,kv.value)
		}
	}
}

func (c *CacheLRU) Add(key string,value Value){
	if ele ,ok := c.cache[key];ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	}else {
		ele := c.ll.PushFront(&entry{key: key,value:value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes{
		c.DeleteOldest()
	}
}

func (c *CacheLRU) Len() int {
	return c.ll.Len()
}
