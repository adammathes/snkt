package web

import (
	"net/http"
)

func Serve(addr string, webroot string) {
	fs := http.FileServer(http.Dir(webroot))
	http.Handle("/", fs)
	http.ListenAndServe(addr, nil)
}
