package hashes

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/**
一致性哈希：
	使用相同的哈希算法计算key和节点名称的哈希值，从而实现一致性哈希
节点变化时：
	根据key值获取到hash，再去获取其节点名称，获取节点名称的逻辑是：
		找到第一个在keys中大于该hash的坐标，再根据该坐标去hashMap中获取节点
*/

// HashFunc 哈希函数，用于计算节点的hash值
type HashFunc func(data []byte) uint32

type Map struct {
	hash     HashFunc       // hash 哈希函数
	replicas int            // replicas 一个真实节点对应虚拟节点的数量
	keys     []int          // keys 表
	hashMap  map[int]string // hashMap key值对应的真实节点名称
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
