package cachefs

import (
	"bytes"
	"io"
	"os"
	"sync"
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

func FuzzHttpCacheFs(f *testing.F) {
	osvalue, err := os.ReadFile("cachefs.go")
	if err != nil {
		f.Fatal(err)
	}
	f.Fuzz(func(t *testing.T, i uint) {
		var wg sync.WaitGroup
		wg.Add(int(i))
		hs := NewHttpCacheFs("./")
		for n := uint(0); n < i; n++ {
			go func() {
				defer wg.Done()
				fs, err := hs.Open("cachefs.go")
				if err != nil {
					t.Fatal(err)
				}
				body, err := io.ReadAll(fs)
				if err != nil && err != io.EOF {
					t.Fatal(err)
				}
				if bytes.Compare(osvalue, body) != 0 {
					t.Fatalf("%s \n!= %s", string(osvalue), string(body))
				}
			}()
		}
		wg.Wait()
	})
}
