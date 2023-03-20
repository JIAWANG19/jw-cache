package dao

import (
	"fmt"
	"log"
	"sync"
)

// Getter 回调函数，但在缓存中获取数据失败时，可以调用回调函数获取数据
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 不同分组，相当于命名空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// Get 根据key获取组内的值，若没有获取到，抛出异常
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[JwCache] hit")
		return v, nil
	}
	// 尝试从其他数据源获取
	return g.load(key)
}

// 从其他数据源获取数据
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// 调用Getter从其他数据源获取数据，若获取到数据，将该数据存入缓存中
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{bytes: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 创建分组
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("空的 Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 根据name获取分组
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}
