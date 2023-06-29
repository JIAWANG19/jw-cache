package lru_cache

import (
	"jw-cache/src/cache"
)

type BaseCacheEvict struct {
	maxBytes  int64                                    // 最大内存
	nowBytes  int64                                    // 当前占用内容
	onEvicted func(key string, value cache.CacheValue) // 缓存淘汰时的回调函数
}

func (baseCache BaseCacheEvict) ableAdd(value cache.CacheValue) bool {
	return baseCache.nowBytes+value.Size() <= baseCache.maxBytes
}

// NowSize 返回缓存占用的大小
func (baseCache BaseCacheEvict) NowSize() int64 {
	return baseCache.nowBytes
}

// MaxCapacity 返回缓存的最大容量
func (baseCache BaseCacheEvict) MaxCapacity() int64 {
	return baseCache.nowBytes
}

// SetMaxCapacity 设置缓存的最大容量
func (baseCache BaseCacheEvict) SetMaxCapacity(maxCapacity int64) error {
	// todo 判断maxCapacity的合理性，根据实际情况进行缓存淘汰
	baseCache.maxBytes = maxCapacity
	return nil
}

// AdjustCapacity 调整缓存的最大容量
func (baseCache BaseCacheEvict) AdjustCapacity(capacity int64) error {
	// todo 判断capacity的合理性，根据实际情况进行缓存淘汰
	baseCache.maxBytes += capacity
	return nil
}
