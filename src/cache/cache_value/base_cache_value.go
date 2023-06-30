package cache_value

import (
	"jw-cache/src/cache"
	"time"
)

type CacheValue interface {
	Size() int64
	Clone() CacheValue
	Refresh() error
	ToString() string
	CustomMethod() interface{}
	SetTimeout(timeout int64)
	SetNeverExpired()
	AdjustTimeout(adjustValue int64)
	Expired() bool
}

type BaseCacheValue struct {
	timeout int64
}

func (value BaseCacheValue) Size() int64 {
	return 0
}

func (value BaseCacheValue) Clone() CacheValue {
	return nil
}

func (value BaseCacheValue) Refresh() error {
	return nil
}

func (value BaseCacheValue) ToString() string {
	return ""
}

func (value BaseCacheValue) CustomMethod() interface{} {
	return nil
}

func (value BaseCacheValue) SetTimeout(timeout int64) {
	value.timeout = timeout
}

func (value BaseCacheValue) SetNeverExpired() {
	value.timeout = cache.NeverTimeout
}

func (value BaseCacheValue) AdjustTimeout(adjustValue int64) {
	value.timeout += adjustValue
}

func (value BaseCacheValue) Expired() bool {
	return value.timeout == cache.NeverTimeout || value.timeout < time.Now().Unix()
}
