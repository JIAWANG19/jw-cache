package cache

import "sync"

type cache struct {
	mu         sync.Mutex
	cache      *Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = New(c.cacheBytes, nil)
	}
	c.cache.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		return
	}
	if v, ok := c.cache.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
