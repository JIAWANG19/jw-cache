package cache_key

type Key struct {
	key string
}

func NewKey(key string) *Key {
	return &Key{key: key}
}

func (key Key) Size() int64 {
	return int64(len(key.key))
}

func (key Key) String() string {
	return key.key
}
