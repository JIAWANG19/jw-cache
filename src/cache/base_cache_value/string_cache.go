package base_cache_value

import (
	"jw-cache/src/cache"
)

type StringCacheValue struct {
	BaseCacheValue
	val     string
	timeout int64
}

func (value StringCacheValue) Size() int64 {
	return int64(len(value.val))
}

func (value StringCacheValue) Clone() cache.CacheValue {
	return StringCacheValue{val: value.val}
}

func (value StringCacheValue) Refresh() error {
	return nil
}

func (value StringCacheValue) ToString() string {
	return value.val
}

func (value StringCacheValue) CustomMethod() interface{} {
	return nil
}
