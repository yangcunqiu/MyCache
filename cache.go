package main

import (
	"MyCache/lru"
	"sync"
)

// cache 并发缓存 把lru.Cache包装起来提供并发特性
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) set(key string, value ByteView) {
	// 加锁
	c.mu.Lock()
	// 释放锁
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Set(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	v, ok := c.lru.Get(key)
	if !ok {
		return
	}
	return v.(ByteView), ok
}
