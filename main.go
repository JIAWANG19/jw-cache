package main

import "jw-cache/src/pgk/log"

// 模拟一个数据源
//var db = map[string]string{
//	"Tom":  "123",
//	"Jack": "234",
//	"Sam":  "345",
//}

// 创建一个组
//func createGroup() *cache.Group {
//	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(
//		func(key string) ([]byte, error) {
//			log.Println("[SlowDB] search key", key)
//			if v, ok := db[key]; ok {
//				return []byte(v), nil
//			}
//			return nil, fmt.Errorf("%s not exist", key)
//		}))
//}

// 根据 addr
//func startCacheServer(addr string, addresses []string, group *cache.Group) {
//	nodes := https.NewHTTPPool(addr)
//	nodes.Set(addresses...)
//	group.RegisterNodes(nodes)
//	log.Println("JWCache is running at", addr)
//	log.Fatal(http.ListenAndServe(addr[7:], nodes))
//}

//func startAPIServer(apiAddr string, group *cache.Group) {
//	http.Handle("/api", http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			key := r.URL.Query().Get("key")
//			view, err := group.Get(key)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusInternalServerError)
//				log.Println(err.Error())
//				return
//			}
//			log.Println(string(view.ByteSlice()))
//			w.Header().Set("Content-Type", "application/octet-stream")
//			w.Write(view.ByteSlice())
//		}))
//	log.Println("source server is running at", apiAddr)
//	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
//}

// 一个开启分布式节点的例子
// 该例子完成了以下操作：
//  1. 开启了三个分布式节点，分别部署在 8001、8002、8003 端口
//  2. 开启了一个入口服务，部署在 9999 端口
//func startAExample() {
//	apiAddr := "http://localhost:9999"
//	addrMap := map[int]string{
//		8001: "http://localhost:8001",
//		8002: "http://localhost:8002",
//		8003: "http://localhost:8003",
//	}
//	var addresses []string
//	for _, v := range addrMap {
//		addresses = append(addresses, v)
//	}
//	for port := range addrMap {
//		port := port
//		group := createGroup()
//		go startCacheServer(addrMap[port], addresses, group)
//	}
//	time.Sleep(2 * time.Second)
//	group := createGroup()
//	startAPIServer(apiAddr, group)
//}

//func info(s string, v interface{}) {
//	logger.Logger.Info("123123", v)
//}

func main() {
	log.Info("%s", "123123")
	log.Info("%s", "123123")
	//info("123")
	//startAExample()
	//var port int
	//var api bool
	//flag.IntVar(&port, "port", 8001, "JWCache server port")
	//flag.BoolVar(&api, "api", false, "Start a api server?")
	//flag.Parse()
	//
	//apiAddr := "http://localhost:9999"
	//addrMap := map[int]string{
	//	8001: "http://localhost:8001",
	//	8002: "http://localhost:8002",
	//	8003: "http://localhost:8003",
	//}
	//
	//var addresses []string
	//for _, v := range addrMap {
	//	addresses = append(addresses, v)
	//}
	//
	//gee := createGroup()
	//if api {
	//	go startAPIServer(apiAddr, gee)
	//}
	//startCacheServer(addrMap[port], addresses, gee)
}
