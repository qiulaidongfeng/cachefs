# cachefs

[English](./README.en.md)

#### 介绍
> 通过go的http.FileServer（http.Dir(path)）可以很方便的创建文件服务器
> go的http.Dir如果没有错误每次Open调用都会进行进行一次os.Open，os.Stat,至少一次的(*os.File).Read调用，这些最终会调用syscall包的函数进行系统调用，即使文件没有修改，这可以优化。
> 本包提供 HttpCacheFs ,可以将 http.Dir(path) 替换为 cachefs.HttpCacheFs(path) ,在被读取文件没有修改（当前通过比较修改时间判断），可以减少系统调用（具体是避免os.Open,(*os.File).Read），提高性能。

#### 实现原理

##### HttpCacheFs
实现了[http.FileSystem](https://pkg.go.dev/net/http#FileSystem) 接口
内部有一个哈希表 key是path value的类型是*CacheFs，用作缓存
如果path已经被哈希表缓存，直接返回缓存，从而避免os.Open

##### CacheFs
实现了[http.File](https://pkg.go.dev/net/http#File)接口
内部用Buf保存文件数据
当Read方法被调用时，先通过比较修改时间判断文件有没有修改
- 如果没有修改，调用Buf的Read方法返回文件数据，避免(*os.File).Read
- 如果有修改，重新读取文件数据，并更新HttpCacheFs的缓存

Close方法为了配合[http.FileServer](https://pkg.go.dev/net/http#FileServer)，永远返回nil,并且不关闭文件句柄

Readdir方法目前返回nil,nil

##### Buf
实现了[io.Reader](https://pkg.go.dev/io#Reader)接口
Buf将[]byte封装成数据流，读取到末尾返回 [io.EOF](https://pkg.go.dev/io#EOF) 后下次读取会从头开始


#### 参与贡献

1.  创建一个issue
2.  Fork 本仓库
3.  新建 Fork_xxx 分支
4.  提交代码
5.  新建 Pull Request

