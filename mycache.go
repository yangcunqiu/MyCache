package main

import (
	"MyCache/lru"
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 命名空间, 具有唯一的名字, 负责和用户交互
type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取源数据的回调
	mainCache cache  // 并发缓存
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func New(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes, lru: lru.New(cacheBytes, nil)},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	value, ok := g.mainCache.get(key)
	if ok {
		log.Printf("[cache] %v hit", key)
		return value, nil
	}
	// 未命中, 调用getter方法
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

// getLocally 获取本地数据
func (g *Group) getLocally(key string) (ByteView, error) {
	// getter回调函数
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	// 拷贝数据
	value := ByteView{b: cloneBytes(bytes)}
	// 回调函数返回的值保存到缓存中
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.set(key, value)
}
