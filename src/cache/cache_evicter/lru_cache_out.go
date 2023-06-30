package cache_evicter

import (
	"container/list"
	k "jw-cache/src/cache/cache_key"
	v "jw-cache/src/cache/cache_value"
	"jw-cache/src/pgk/log"
)

// LRUCacheEvict LRU 缓存淘汰策略
type LRUCacheEvict struct {
	BaseCacheEvicter
	evictList *list.List
	cache     map[*k.Key]*list.Element
	//mu         sync.Mutex
}

type cacheNode struct {
	key   *k.Key
	value v.CacheValue
}

func NewLRUCache(maxBytes int64, onEvicted func(*k.Key, v.CacheValue)) *LRUCacheEvict {
	return &LRUCacheEvict{
		BaseCacheEvicter{maxBytes: maxBytes, nowBytes: 0, onEvicted: onEvicted},
		list.New(),
		make(map[*k.Key]*list.Element)}
}

func (cache LRUCacheEvict) Add(key *k.Key, value v.CacheValue) {
	if !cache.BeforeAdd(key, value) {
		return
	}
	if _, exist := cache.cache[key]; exist {
		log.Debug("[Cache] 缓存值已存在, key: %s", key)
		return
	}
	cache.nowBytes += key.Size() + value.Size()
	ele := cache.evictList.PushFront(&cacheNode{key: key, value: value})
	cache.cache[key] = ele
	log.Debug("[Cache] 缓存值添加成功: {key: %s, value: %s}", key.String(), value.ToString())
}

func (cache LRUCacheEvict) Get(key *k.Key) (value v.CacheValue, exist bool) {
	if val, ok := cache.cache[key]; ok {
		cache.evictList.MoveToFront(val)
		kv := val.Value.(*cacheNode)
		return kv.value, true
	}
	return nil, false
}

func (cache LRUCacheEvict) Update(key *k.Key, value v.CacheValue) error {
	//if val , exist := cache.cache[key]; exist {
	//} else {
	//
	//}
	return nil
}

func (cache LRUCacheEvict) Delete(key *k.Key) error {
	if ele, ok := cache.cache[key]; ok {
		cache.evictList.Remove(ele)
		delete(cache.cache, key)
		kv := ele.Value.(*cacheNode)
		cache.nowBytes -= key.Size() + kv.value.Size()
		log.Debug("[Cache] 缓存删除: %v", kv)
	} else {
		log.Debug("[Cache] 缓存值不存在, key: %s", key)
	}
	return nil
}

func (cache LRUCacheEvict) Clear() error {
	return nil
}

func (cache LRUCacheEvict) Keys() []*k.Key {
	return nil
}

func (cache LRUCacheEvict) Has(key *k.Key) bool {
	return true
}

func (cache LRUCacheEvict) Evict() error {
	return nil
}
