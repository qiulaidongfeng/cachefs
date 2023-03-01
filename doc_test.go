package cachefs

import (
	"net/http"
)

func ExampleHttpCacheFs() {
	path := "path"
	http.Handle("/", http.FileServer(NewHttpCacheFs(path)))
}
