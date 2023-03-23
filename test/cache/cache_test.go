package cache

import (
	"JWCache/cache"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	cache := cache.New(1024, nil)
	cache.Add("key1", String("value1"))
	if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "value1" {
		t.Fatalf("缓存失败")
	}
	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("缓存失败")
	}
}

func TestRemove(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	bytes := len(k1 + k2 + v1 + v2)
	cache := cache.New(int64(bytes), nil)
	cache.Add(k1, String(v1))
	cache.Add(k2, String(v2))
	cache.Add(k3, String(v3))
	if _, ok := cache.Get("key1"); ok || cache.Len() != 2 {
		t.Fatalf("缓存失败")
	}
}
