package hashes

import (
	"JWCache/hashes"
	"strconv"
	"testing"
)

func TestHashes(t *testing.T) {
	hash := hashes.New(3, func(key []byte) uint32 {
		// 直接返回传入字符串对应的数字
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// 应当生产以下虚拟节点
	// 2 4 6 12 14 16 22 24 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	var startTest = func() {
		for k, v := range testCases {
			if hash.Get(k) != v {
				t.Errorf("该 %s 对应的value值应当是 %s", k, v)
			}
		}
	}

	startTest()
	// 添加一个节点
	hash.Add("8")
	testCases["27"] = "8"
	startTest()
}
