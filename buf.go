package cachefs

import (
	"io"
	"sync/atomic"
)

var _ io.Reader = (*Buf)(nil)

// Buf 循环缓存
// 当到末尾是不会返回 [io.EOF] ,会下次从头读取
type Buf struct {
	i      atomic.Int32
	buflen int32
	buf    []byte
}

// NewBuf 创建循环缓存
func NewBuf(buf []byte) *Buf {
	ret := &Buf{buf: buf}
	ret.buflen = int32(len(buf))
	return ret
}

// Read 实现io.Reader接口,读取缓存内容
// 返回 [io.EOF] 后下次读取会从头开始
func (buf *Buf) Read(p []byte) (n int, err error) {
	if buf.Empty() {
		if len(p) == 0 {
			return 0, nil
		}
		buf.i.Store(0)
		return 0, io.EOF
	}
	lp := int32(len(p))
	ni := buf.i.Add(lp)
	n = copy(p, buf.buf[ni-lp:])
	return n, nil
}

// Seek 实现io.Seeker接口
func (buf *Buf) Seek(offset int64, whence int) (int64, error) {
	off := int32(offset)
	switch whence {
	case io.SeekStart:
		buf.i.Store(off)
	case io.SeekCurrent:
		buf.i.Add(off)
	case io.SeekEnd:
		buf.i.Store(buf.buflen - off)
	}
	return int64(buf.i.Load()), nil
}

// Empty 判断缓存是否到末尾
func (buf *Buf) Empty() bool {
	return buf.i.Load() >= buf.buflen
}
