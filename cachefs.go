// Package cachefs 实现了缓存文件系统
//
// 缓存文件系统用于在 [http.FileSystem] 默认实现需要优化
// 通过比较修改时间，缓存文件系统在Open和Read时如果没有修改，减少系统调用（具体是避免os.Open,(*os.File).Read）
package cachefs

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HttpCacheFs Http缓存文件系统
type HttpCacheFs struct {
	path string
	fd   sync.Map //key=string value=*CacheFs
}

// NewHttpCacheFs 创建NewHttpCacheFs
//   - path 是相对的路径 相当于 [http.Dir]
func NewHttpCacheFs(path string) *HttpCacheFs {
	return &HttpCacheFs{path: path}
}

// Open 实现 [http.FileSystem] 接口
func (fs *HttpCacheFs) Open(name string) (http.File, error) {
	fdname := filepath.Join(fs.path, name)
	value, ok := fs.fd.Load(fdname)
	if ok { //如果已被缓存
		return value.(*CacheFs).Copy(), nil
	}
	cache, err := NewCacheFs(fdname)
	if err != nil {
		return nil, err
	}
	cache.fs = fs
	fs.fd.Store(fdname, cache)
	return cache.Copy(), nil
}

// CacheFs 缓存文件系统
//
// 通过比较修改时间，在Open和Read时如果没有修改，减少系统调用（具体是避免os.Open,(*os.File).Read）
type CacheFs struct {
	fs      *HttpCacheFs
	fd      *os.File
	fdinfo  os.FileInfo
	modtime time.Time //缓存创建时修改时间
	buf     Buf
}

// NewCacheFs 创建缓存文件系统
func NewCacheFs(name string) (fs *CacheFs, err error) {
	fd, fdinfo, err := openAndStat(name)
	if err != nil {
		return nil, err
	}
	modtime := fdinfo.ModTime() //获取修改时间
	file, err := io.ReadAll(fd) //读取全部内容
	if err != nil {
		return nil, err
	}
	ret := &CacheFs{fd: fd, modtime: modtime, buf: NewBuf(file)} //创建缓存文件系统，文件全部内容被缓存
	return ret, nil
}

// openAndStat 调用 [os.Open] 和 [os.File.Stat]
func openAndStat(name string) (*os.File, fs.FileInfo, error) {
	fd, err := os.Open(name) //打开
	if err != nil {
		return nil, nil, err
	}
	fdinfo, err := fd.Stat()
	if err != nil {
		return nil, nil, err
	}
	return fd, fdinfo, nil
}

// resetRead 重新读取
func (fs *CacheFs) resetRead() error {
	fd, fdinfo, err := openAndStat(fs.fd.Name())
	if err != nil {
		return err
	}
	fs.fd = fd
	fs.modtime = fdinfo.ModTime()
	file, err := io.ReadAll(fd) //读取全部内容
	if err != nil {
		return err
	}
	fs.buf = NewBuf(file)
	fs.fs.fd.Store(fs.fd.Name(), fs.Copy())
	return nil
}

// Read 实现 [io.Reader] 接口，如果没有修改，将返回缓存内容
func (fs *CacheFs) Read(p []byte) (n int, err error) {
	ok, err := fs.isNoRevise()
	if err != nil {
		return 0, err
	}
	if ok { //通过比较缓存时修改时间，判断是否修改，没有直接从缓存读取
		return fs.buf.Read(p)
	}
	//有修改，重新缓存再读
	err = fs.resetRead()
	if err != nil {
		return 0, err
	}
	return fs.Read(p)
}

// Seek 实现 [io.Seeker] 接口，如果没有修改，将移动缓存内容偏移量
func (fs *CacheFs) Seek(offset int64, whence int) (int64, error) {
	ok, err := fs.isNoRevise()
	if err != nil {
		return 0, err
	}
	if ok { //通过比较缓存时修改时间，判断是否修改，没有直接移动缓存内容偏移量
		return fs.buf.Seek(offset, whence)
	}
	//有修改，重新缓存再移动缓存内容偏移量
	err = fs.resetRead()
	if err != nil {
		return 0, err
	}
	return fs.Seek(offset, whence)
}

// Readdir 返回目录信息
func (fs *CacheFs) Readdir(count int) ([]fs.FileInfo, error) {
	ok, err := fs.isNoRevise()
	if err != nil {
		return nil, err
	}
	if ok { //如果文件没有修改
		return fs.fd.Readdir(count)
	}
	//如果文件有修改
	fs.resetRead()
	return fs.Readdir(count)
}

// isNoRevise 返回系统文件是否没有修改
func (fs *CacheFs) isNoRevise() (bool, error) {
	var fdinfo os.FileInfo
	if fs.fdinfo != nil {
		fdinfo = fs.fdinfo
	} else {
		var err error
		fdinfo, err = fs.fd.Stat()
		if err != nil {
			return false, err
		}
	}
	nowtime := fdinfo.ModTime()           //获取文件现在修改时间
	return fs.modtime.Equal(nowtime), nil //通过比较缓存时修改时间，判断是否没有修改
}

// Close 关闭
//
// 为了能配合 [http.FileServer] ，永远返回nil，并且不关闭文件句柄
func (fs *CacheFs) Close() error {
	return nil
}

// Stat 返回文件信息
//
// 如果没有修改，避免 [os.Open] 系统调用
func (fs *CacheFs) Stat() (fs.FileInfo, error) {
	ok, err := fs.isNoRevise()
	if err != nil {
		return nil, err
	}
	if ok { //通过比较缓存时修改时间，判断是否修改，没有直接从缓存读取
		fs.fdinfo, err = fs.fd.Stat()
		return fs.fdinfo, err
	}
	fs.resetRead()
	return fs.Stat()
}

// Copy 对自身进行浅拷贝
func (fs *CacheFs) Copy() *CacheFs {
	ret := *fs
	ret.buf = ret.buf.Copy()
	return &ret
}
