package cachefs

import (
	"net/http"
)

func ExampleHttpCacheFs() {
	path := "path"
	http.Handle("/", http.FileServer(NewHttpCacheFs(path)))
}

func ExampleHttpCacheFs_onefile() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "path"
		fs := NewHttpCacheFs(path)
		f, err := fs.Open("index.html")
		if err != nil {
			panic(err)
		}
		i, err := f.Stat()
		if err != nil {
			panic(err)
		}
		http.ServeContent(w, r, "index.html", i.ModTime(), f)
	})
}
