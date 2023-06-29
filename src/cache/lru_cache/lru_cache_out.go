package lru_cache

import (
	"container/list"
	"jw-cache/src/cache"
)

type LRUCacheEvict struct {
	BaseCacheEvict
	evictList *list.List
	cache     map[string]*list.Element
	//mu         sync.Mutex
}

type cacheNode struct {
	key   string
	value cache.CacheValue
}

func (cache LRUCacheEvict) Add(key string, value cache.CacheValue) {
	if _, exist := cache.cache[key]; exist {
	}
}

func (cache LRUCacheEvict) Get(key string) (value cache.CacheValue, exist bool) {
	if val, ok := cache.cache[key]; ok {
		cache.evictList.MoveToFront(val)
		kv := val.Value.(*cacheNode)
		return kv.value, true
	}
	return nil, false
}

func (cache LRUCacheEvict) Update(key string, value cache.CacheValue) error {
	//if val, ok := cache.cache[key]; ok {
	//
	//}
	return nil
}

func (cache LRUCacheEvict) Delete(key string) error {
	return nil
}

func (cache LRUCacheEvict) Clear() error {
	return nil
}

func (cache LRUCacheEvict) Keys() []string {
	return nil
}

func (cache LRUCacheEvict) Has(key string) bool {
	return true
}

func (cache LRUCacheEvict) Evict() error {
	return nil
}
