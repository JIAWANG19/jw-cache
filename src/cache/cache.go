package cache

import (
	"container/list"
)

// Cache 定义了LRU Cache的基本数据结构
type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nowBytes  int64                         // 当前已使用的内存
	ll        *list.List                    // 双向链表，记录访问频率，越是最近被使用的记录越不容易被淘汰
	cache     map[string]*list.Element      // 哈希表，记录每个键对应的值在链表中的位置
	OnEvicted func(key string, value Value) // 记录被删除时的回调函数，可选参数
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

// CacheEvicter 缓存淘汰接口
type CacheEvicter interface {
	Add(key string, value CacheValue)              // Add 添加缓存值
	Get(key string) (value CacheValue, exist bool) // Get 根据键获取缓存值，并返回是否存在
	Update(key string, value CacheValue) error     // Update 修改缓存值
	Delete(key string) error                       // Delete 删除指定键的缓存值
	Clear() error                                  // Clear 清除所有缓存值
	Keys() []string                                // Keys 返回缓存中的所有键
	Has(key string) bool                           // Has 检查指定键是否存在于缓存中
	Evict() error                                  // Evict 根据一定的策略驱逐缓存值

	NowSize() int64                         // Size 返回缓存的总大小（以字节为单位）
	MaxCapacity() int64                     // MaxCapacity 返回缓存的最大容量
	AdjustCapacity(capacity int64) error    // AdjustCapacity 调整缓存容量大小
	SetMaxCapacity(maxCapacity int64) error // SetMaxCapacity 设置缓存的最大容量
}

const NeverTimeout = -1

// CacheValue 缓存值接口
type CacheValue interface {

	// Size 返回缓存值的大小（以字节为单位）
	Size() int64

	// Clone 返回缓存值的副本
	Clone() CacheValue

	// Refresh 刷新缓存值
	Refresh() error

	// ToString 将缓存值转换为字符串
	ToString() string

	// CustomMethod 自定义方法，根据需要实现
	CustomMethod() interface{}

	// SetTimeout 过期时间
	SetTimeout(timeout int64)

	// SetNeverExpired 设置永不过期
	SetNeverExpired()

	// AdjustTimeout 调整过期时间
	AdjustTimeout(adjustValue int64)

	// Expired 检查缓存值是否已过期
	Expired() bool
}
