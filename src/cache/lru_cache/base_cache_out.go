package lru_cache

import (
	"jw-cache/src/cache"
)

type BaseCacheEvict struct {
	maxBytes  int64                                    // 最大内存
	nowBytes  int64                                    // 当前占用内容
	onEvicted func(key string, value cache.CacheValue) // 缓存淘汰时的回调函数
}

func (baseCache BaseCacheEvict) NowSize() int64 {
	return baseCache.nowBytes
}

func (baseCache BaseCacheEvict) MaxCapacity() int64 {
	return baseCache.nowBytes
}
func (baseCache BaseCacheEvict) SetMaxCapacity(maxCapacity int64) error {
	// todo 判断maxCapacity的合理性，根据实际情况进行缓存淘汰
	baseCache.maxBytes = maxCapacity
	return nil
}
func (baseCache BaseCacheEvict) AdjustCapacity(capacity int64) error {
	// todo 判断capacity的合理性，根据实际情况进行缓存淘汰
	baseCache.maxBytes += capacity
	return nil
}
