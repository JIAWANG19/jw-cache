package lru_cache

import (
	"container/list"
	"jw-cache/src/cache"
	"jw-cache/src/pgk/log"
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
	if !cache.ableAdd(value) {
		log.Debug("[Cache] 缓存已满, 现容量: %s, 最大容量: %d", cache.nowBytes, cache.maxBytes)
	}
	if _, exist := cache.cache[key]; exist {
		log.Debug("[Cache] 缓存值已存在, key: %s", key)
	}
	cache.nowBytes += int64(len(key)) + value.Size()
	ele := cache.evictList.PushFront(&cacheNode{key: key, value: value})
	cache.cache[key] = ele
	log.Debug("[Cache] 缓存值添加成功: %v", ele)
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
	//return nil
	return nil
}

func (cache LRUCacheEvict) Delete(key string) error {
	if ele, ok := cache.cache[key]; ok {
		cache.evictList.Remove(ele)
		delete(cache.cache, key)
		kv := ele.Value.(*cacheNode)
		cache.nowBytes -= int64(len(key)) + kv.value.Size()
		log.Debug("[Cache] 缓存删除: %v", kv)
	} else {
		log.Debug("[Cache] 缓存值不存在, key: %s", key)
	}
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
