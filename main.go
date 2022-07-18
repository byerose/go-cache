package main

import (
	"log"
	"net/http"
)

type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!"))
}

func main() {
	var s server
	//第一个参数是服务启动的地址，第二个参数是 Handler，任何实现了 ServeHTTP 方法的对象都可以作为 HTTP 的 Handler。
	http.ListenAndServe("localhost:9999", &s)
}
