package jwCache

import (
	"JWCache/jwCache"
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f jwCache.Getter = jwCache.GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("回调函数失败")
	}
}

// 模拟一个数据源
var db = map[string]string{
	"Tom":  "123",
	"Jack": "234",
	"Sam":  "345",
}

func TestGroupGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	cache := jwCache.NewGroup("mysql", 2<<10, jwCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	for k, v := range db {
		if view, err := cache.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s\n", k)
		}
		if _, err := cache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := cache.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be emtry, but %s got", view)
	}
}
