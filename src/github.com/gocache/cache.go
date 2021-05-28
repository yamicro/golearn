package gocache

import (
	"sync"
)

type Cache struct{
	mu sync.Mutex
	lru *CacheLRU
	CacheBytes int64 //cache中储存的值
}

func (c *Cache) Add(key string,value ByteView)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		c.lru = NewCache(c.CacheBytes,nil) // lazy initialzation
	}
	c.lru.Add(key,value)
}

func (c *Cache) Get(key string) (value ByteView,ok bool){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}