package hashes

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/**
一致性哈希是一种用于分布式系统中数据存储和负载均衡的算法。在一致性哈希中，每个节点都被映射到一个哈希环上，数据的key也被映射到哈希环上。
通过使用相同的哈希算法，节点名称和key都可以被转换为哈希值，从而可以将数据key映射到最接近的节点上。
当有节点发生变化时，只需找到比该节点哈希值大的第一个节点，即可将数据迁移到该节点上。
这种方法避免了大规模数据移动的问题，同时保持了数据的分布均衡。
*/

// HashFunc 哈希函数，用于计算的哈希值
type HashFunc func(data []byte) uint32

type Map struct {
	hash     HashFunc       // hash 用于计算哈希值的哈希函数
	replicas int            // replicas 一个真实节点对应虚拟节点的数量
	keys     []int          // keys 该变量是一个有序列表，包含所有的虚拟节点
	hashMap  map[int]string // hashMap 该变量是一个哈希表，存储虚拟节点的哈希值和对应的真实节点名称
}

// New 创建一个 Map，如果哈希函数为空，使用默认的哈希函数
func New(replicas int, hashFunc HashFunc) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hashFunc,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加一个0个或多个节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 根据哈希函数获取虚拟节点的哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 将虚拟节点的哈希值添加到
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 根据key获取节点名称
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
