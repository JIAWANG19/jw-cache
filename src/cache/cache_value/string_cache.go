package cache_value

type StringCacheValue struct {
	BaseCacheValue
	val string
}

func NewStringValue(value string, timeout int64) *StringCacheValue {
	return &StringCacheValue{
		BaseCacheValue{timeout: timeout},
		value,
	}
}

func (value StringCacheValue) Size() int64 {
	return int64(len(value.val))
}

func (value StringCacheValue) Clone() CacheValue {
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
