package dao

import (
	"JWCache/nodes"
	"JWCache/singleflight"
	"fmt"
	"log"
	"sync"
)

// Getter 回调函数，但在缓存中获取数据失败时，可以调用回调函数获取数据
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 不同分组，相当于命名空间
type Group struct {
	name      string              // 组名
	getter    Getter              // 回调函数，当在缓存中没有查询到数据时，去调用回调函数获取数据
	mainCache cache               // 缓存的具体实现
	nodes     nodes.NodePicker    // 节点
	loader    *singleflight.Group // 防止缓存击穿的实现
}

// RegisterNodes 注册节点
func (g *Group) RegisterNodes(nodes nodes.NodePicker) {
	if g.nodes != nil {
		panic("RegisterNodePicker called more than once")
	}
	g.nodes = nodes
}

func (g *Group) load(key string) (value ByteView, err error) {
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.nodes != nil {
			if node, ok := g.nodes.PickNode(key); ok {
				if value, err = g.GetFromNode(node, key); err == nil {
					return value, nil
				}
				log.Println("[JWCache] Failed to get for node", err)
			}
		}
		return g.getLocally(key)
	})

	if err == nil {
		return view.(ByteView), nil
	}
	return
	//if g.nodes != nil {
	//	if node, ok := g.nodes.PickNode(key); ok {
	//		if value, err = g.GetFromNode(node, key); err == nil {
	//			return value, nil
	//		}
	//		log.Println("[JWCache] Failed to get for node", err)
	//	}
	//}
	//return g.getLocally(key)
}

// GetFromNode 从指定节点中获取数据
func (g *Group) GetFromNode(node nodes.NodeGetter, key string) (ByteView, error) {
	bytes, err := node.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{bytes: bytes}, nil
}

// Get 根据key获取组内的值，若没有获取到，抛出异常
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		// todo 为了方便测试，这里暂时不打印该日志
		//log.Println("[JwCache] hit")
		return v, nil
	}
	// 尝试从其他数据源获取
	return g.load(key)
}

// 从其他数据源获取数据
//func (g *Group) load(key string) (value ByteView, err error) {
//	return g.getLocally(key)
//}

// 调用Getter从其他数据源获取数据，若获取到数据，将该数据存入缓存中
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{bytes: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 创建分组
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("空的 Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup 根据name获取分组
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}
