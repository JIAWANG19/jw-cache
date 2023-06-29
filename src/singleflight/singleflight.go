package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu      sync.Mutex
	callMap map[string]*call
}

// Do 防止缓存击穿的实现，当相同的key并发的请求时，该方法可以保证fn函数只被调用一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()           // 先上锁
	if g.callMap == nil { // 延迟加载
		g.callMap = make(map[string]*call)
	}
	if c, ok := g.callMap[key]; ok { // 如果有相同的key正在请求，则等待
		g.mu.Unlock()       // 解锁
		c.wg.Wait()         // 等待key请求完成
		return c.val, c.err // 返回key请求的结果
	}
	aCall := new(call)
	aCall.wg.Add(1)        // 发起请求前加锁，使请求结束前的所有与该请求相同的key阻塞
	g.callMap[key] = aCall // 添加到 g.callMap
	g.mu.Unlock()

	aCall.val, aCall.err = fn() // 调用方法获取key的值
	aCall.wg.Done()             // 请求解锁

	g.mu.Lock()
	delete(g.callMap, key) // 更新 g.callMap
	g.mu.Unlock()

	return aCall.val, aCall.err
}
