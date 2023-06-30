package cache_evicter

import (
	k "jw-cache/src/cache/cache_key"
	v "jw-cache/src/cache/cache_value"
	"jw-cache/src/pgk/log"
	"strconv"
)

// CacheEvicter 缓存淘汰接口
type CacheEvicter interface {
	Add(key *k.Key, value v.CacheValue)              // Add 添加缓存值
	Get(key *k.Key) (value v.CacheValue, exist bool) // Get 根据键获取缓存值，并返回是否存在
	Update(key *k.Key, value v.CacheValue) error     // Update 修改缓存值
	Delete(key *k.Key) error                         // Delete 删除指定键的缓存值
	Clear() error                                    // Clear 清除所有缓存值
	Keys() []*k.Key                                  // Keys 返回缓存中的所有键
	Has(key *k.Key) bool                             // Has 检查指定键是否存在于缓存中
	Evict() error                                    // Evict 根据一定的策略驱逐缓存值

	NowSize() int64                         // Size 返回缓存的总大小（以字节为单位）
	MaxCapacity() int64                     // MaxCapacity 返回缓存的最大容量
	AdjustCapacity(capacity int64) error    // AdjustCapacity 调整缓存容量大小
	SetMaxCapacity(maxCapacity int64) error // SetMaxCapacity 设置缓存的最大容量
}

type CacheBefore interface {
	BeforeAdd(key *k.Key, value v.CacheValue) bool    // BeforeAdd 添加前调用
	BeforeUpdate(key *k.Key, value v.CacheValue) bool // BeforeUpdate 修改前调用
}

type CacheAfter interface {
}

type BaseCacheEvicter struct {
	maxBytes  int64                                // 最大内存
	nowBytes  int64                                // 当前占用内容
	onEvicted func(key *k.Key, value v.CacheValue) // 缓存淘汰时的回调函数
}

func (cache *BaseCacheEvicter) Add(key *k.Key, value v.CacheValue) {
}

func (cache *BaseCacheEvicter) BeforeAdd(key *k.Key, value v.CacheValue) bool {
	if value == nil {
		log.Debug("[Cache] value is null, key: %s, value: %v", key, value)
		return false
	}
	if cache.nowBytes+key.Size()+value.Size() > cache.maxBytes {
		log.Debug("[Cache] 缓存已满, 现容量: %s, 最大容量: %d", strconv.FormatInt(cache.nowBytes, 10), cache.maxBytes)
		return false
	}
	return true
}

func (cache *BaseCacheEvicter) BeforeUpdate(key *k.Key, value v.CacheValue) bool {
	return false
}

func (cache *BaseCacheEvicter) Get(key *k.Key) (value v.CacheValue, exist bool) {
	return nil, false
}

func (cache *BaseCacheEvicter) Update(key *k.Key, value v.CacheValue) error {
	return nil
}

func (cache *BaseCacheEvicter) Delete(key *k.Key) error {
	return nil
}

func (cache *BaseCacheEvicter) Clear() error {
	return nil
}

func (cache *BaseCacheEvicter) Keys() []*k.Key {
	return nil
}

func (cache *BaseCacheEvicter) Has(key *k.Key) bool {
	return false
}

func (cache *BaseCacheEvicter) Evict() error {
	return nil
}

// NowSize 返回缓存占用的大小
func (cache *BaseCacheEvicter) NowSize() int64 {
	return cache.nowBytes
}

// MaxCapacity 返回缓存的最大容量
func (cache *BaseCacheEvicter) MaxCapacity() int64 {
	return cache.nowBytes
}

// SetMaxCapacity 设置缓存的最大容量
func (cache *BaseCacheEvicter) SetMaxCapacity(maxCapacity int64) error {
	// todo 判断maxCapacity的合理性，根据实际情况进行缓存淘汰
	cache.maxBytes = maxCapacity
	return nil
}

// AdjustCapacity 调整缓存的最大容量
func (cache *BaseCacheEvicter) AdjustCapacity(capacity int64) error {
	// todo 判断capacity的合理性，根据实际情况进行缓存淘汰
	cache.maxBytes += capacity
	return nil
}
