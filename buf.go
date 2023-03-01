package cachefs

import (
	"io"
)

var _ io.Reader = (*Buf)(nil)

// Buf 循环缓存
type Buf struct {
	buflen int
	i      int
	buf    []byte
}

// NewBuf 创建循环缓存
func NewBuf(buf []byte) Buf {
	return Buf{buf: buf, buflen: len(buf)}
}

// Read 实现 [io.Reader] 接口,读取缓存内容
// 读取到末尾返回 [io.EOF] 后下次读取会从头开始
func (buf *Buf) Read(p []byte) (n int, err error) {
	if buf.Empty() {
		if len(p) == 0 {
			return 0, nil
		}
		buf.i = 0
		return 0, io.EOF
	}
	n = copy(p, buf.buf[buf.i:])
	buf.i += n
	return n, nil
}

// Seek 实现 [io.Seeker] 接口
func (buf *Buf) Seek(offset int64, whence int) (int64, error) {
	off := int(offset)
	switch whence {
	case io.SeekStart:
		buf.i = off
	case io.SeekCurrent:
		buf.i += off
	case io.SeekEnd:
		buf.i = buf.buflen - off
	}
	return int64(buf.i), nil
}

// Empty 判断缓存是否到末尾
func (buf *Buf) Empty() bool {
	return buf.i >= buf.buflen
}

// Copy 浅拷贝一个循环缓存
func (buf *Buf) Copy() Buf {
	b := *buf
	return b
}
