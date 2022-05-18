package main

import (
	"fmt"
	"net/http"
)

type Engine struct{}

func (*Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL path = %s \n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "head[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 was not found %s \n", r.URL.Path)
	}
}

func main() {
	var engine *Engine
	http.HandleFunc("/nihao", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "nihao %s \n", r.URL.Path)
	})

	//监听的func没使用默认的，上面的设置handlefunc的操作会被无视
	http.ListenAndServe(":8080", engine)
}
