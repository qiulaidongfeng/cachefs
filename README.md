# cachefs

#### 介绍
> go的http.Dir每次读取都会进行系统调用，即使文件没有修改，这可以优化。
> 本包提供 HttpCacheFs ,可以将 http.Dir(path) 替换为 cachefs.HttpCacheFs(path) ,在被读取文件没有修改（当前通过比较修改时间判断），可以避免系统调用，提高性能。


#### 参与贡献

1.  创建一个issue
2.  Fork 本仓库
3.  新建 Fork_xxx 分支
4.  提交代码
5.  新建 Pull Request

