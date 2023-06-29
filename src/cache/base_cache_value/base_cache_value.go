package base_cache_value

import (
	"jw-cache/src/cache"
	"time"
)

type BaseCacheValue struct {
	timeout int64
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
