package main

import (
	"JWCache/controller"
	"JWCache/jwCache"
	"fmt"
	"log"
	"net/http"
)

// 模拟一个数据源
var db = map[string]string{
	"Tom":  "123",
	"Jack": "234",
	"Sam":  "345",
}

func main() {
	jwCache.NewGroup("scores", 2<<10, jwCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := controller.NewPool(addr)
	log.Println("jwCache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
