package jwCache

// ByteView 只读数据结构，用于支持并发操作
type ByteView struct {
	bytes []byte
}

// Len 继承 Value 接口的方法
func (v ByteView) Len() int {
	return len(v.bytes)
}

// ByteSlice 返回当前数据的拷贝
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.bytes)
}

func (v ByteView) String() string {
	return string(v.bytes)
}

// 拷贝数据
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
