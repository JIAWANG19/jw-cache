package list

import (
	"jw-cache/src/cache"
)

type node struct {
	val  cache.CacheValue
	prev *node
	next *node
}

// LinkedList 双向链表结构
type LinkedList struct {
	size     int
	nowBytes int64
}

//func (list *LinkedList) MoveToFront() error {
//
//}
