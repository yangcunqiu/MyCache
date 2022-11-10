package main

// ByteView 抽象了一个只读数据结构 ByteView 用来表示缓存值
type ByteView struct {
	// 缓存的真实值
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回b的拷贝, 防止缓存被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c

}
