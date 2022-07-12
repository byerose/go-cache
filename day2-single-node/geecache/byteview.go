package geecache

// A ByteView holds an immutable view of bytes.
//使用[]byte保存真实的缓存值，可以存储字符串图片等类型
type ByteView struct {
	b []byte
}

// Len returns the view's length
//计算缓存占用内存
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
//为什么不直接使用结构体赋值（深拷贝），而是通过内置函数copy()
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
//转换为字符串
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
