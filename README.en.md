# cachefs

[English](./README.en.md)

#### Introduction
>Through http. FileServer (http. Dir (path)) of go, you can easily create a file server
>If there is no error in http.Dir of go, each Open call will make one os.Open, os.Stat, and at least one (* os.File). Read call. These will eventually call the syscall package function for system calls. Even if the file is not modified, this can be optimized.
>This package provides HttpCacheFs, which can replace http.Dir (path) with cachefs.HttpCacheFs (path). When the read file is not modified (currently judged by comparing the modification time), it can reduce system calls (specifically avoid os.Open, (* os. File). Read) and improve performance.

#### Implementation principle

##### HttpCacheFs

Implemented [http. FileSystem](https://pkg.go.dev/net/http#FileSystem) Interface

There is an internal hash table whose key is path value and whose type is *CacheFs, which is used as cache

If the path has been cached by the hash table, the cache is returned directly to avoid os.Open

##### CacheFs

Implemented [http. File](https://pkg.go.dev/net/http#File) Interface

Use Buf internally to save file data and keep the (*os.File) handle of the file

When the Read or Seek method is called, first determine whether the file has been modified by comparing the modification time

- If not modified, call Buf's Read or Seek method to avoid (*os.File). Read or (*os.File). Seek

- If there is any modification, re-read the file data, update the cache of HttpCacheFs, and then call the Read or Seek method

The Close method is used to match [http. FileServer](https://pkg.go.dev/net/http#FileServer), always return nil, and do not close the file handle

When the Readdir or Stat method is called, first determine whether the file has been modified by comparing the modification time

- If it is not modified, call Readdir or Stat of os.File

- If there is any modification, re-read the file data, update the cache of HttpCacheFs, and then call the Readdir or Stat method

##### Buf

Implemented [io. Reader](https://pkg.go.dev/io#Reader) Interface

Implemented [io. Seeker](https://pkg.go.dev/io#Seeker) Interface

Buf encapsulates []byte into a data stream, reads it to the end and returns [io. EOF](https://pkg.go.dev/io#EOF) The next reading will start from the beginning

#### Participation contribution
1. Create an issue
2. Fork warehouse
3. Create a new Fork_ Xxx branch
4. Submit code
5. Create a new Pull Request