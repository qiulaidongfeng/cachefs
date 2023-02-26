package cachefs

import (
	"io"
	"testing"
)

func TestHttpCacheFs(t *testing.T) {
	hfs := NewHttpCacheFs("./")
	fs1, err := hfs.Open("cachefs.go")
	if err != nil {
		t.Fatal(err)
	}
	defer fs1.Close()
	fs2, err := hfs.Open("cachefs.go")
	if err != nil {
		t.Fatal(err)
	}
	fs2.Close()
}

func TestCacheFsRead(t *testing.T) {
	hfs := NewHttpCacheFs("./")
	fs1, err := hfs.Open("cachefs.go")
	if err != nil {
		t.Fatal(err)
	}
	defer fs1.Close()
	p := make([]byte, 10240)
	_, err = fs1.Read(p)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	_, err = fs1.Read(p)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	_, err = fs1.Read(p)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	_, err = fs1.Read(p)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	_, err = fs1.Read(p)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	_, err = fs1.Read(p)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
}
