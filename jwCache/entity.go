package jwCache

import "container/list"

// Cache 使用LRU(最近最少未使用)作为淘汰策略，于是需要ll记录访问频率
type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nowBytes  int64                         // 已使用的内存
	ll        *list.List                    // 记录访问频率，越是最近被使用的记录越不容易被淘汰
	cache     map[string]*list.Element      // map
	OnEvicted func(key string, value Value) // 记录被删除时的回调函数
}

type entry struct {
	key   string
	value Value
}

// Value 所有的值都需要实现这个接口
type Value interface {
	Len() int //返回值占用的内存大小
}

// Get 查找功能, @value 返回值，@success 是查找成功
// 若查询成功，将该节点移至队首，表示最近被使用
func (c *Cache) Get(key string) (value Value, success bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// Remove 删除(缓存淘汰)
func (c *Cache) Remove() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增/修改 操作结束后将被操作的节点移到队列的队首，若发现已使用内存超过了最大内存，则调用回收方法
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.Remove()
	}
}

// Len 方便测试
func (c *Cache) Len() int {
	return c.ll.Len()
}

// New 实例化 Cache
func New(maxBytes int64, OnEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}
