package main

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

// key v, value v (1 obj)
// var m = map[[12]byte]int{}

// func init() {
// 	for i := 0; i < 1000000; i++ {
// 		var key [12]byte
// 		copy(key[:], strconv.Itoa(i))
// 		m[key] = i
// 	}
// }

// key ptr, value v (n obj)
// var mstr = map[string]int{}

// func init() {
// 	for i := 0; i < 1000000; i++ {
// 		var key string
// 		key = strconv.Itoa(i)
// 		mstr[key] = i
// 	}
// }

// key ptr, value ptr (n obj)
var mstrptr = map[string]*int{}

func init() {
	for i := 0; i < 1000000; i++ {
		var key string
		t := i
		key = strconv.Itoa(i)
		mstrptr[key] = &t
	}
}

// go http server.
func main() {
	// 编写一个标准的http服务
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// 测试复用内存的cpu、mem情况
func pprofWhenReUseSlice() { return }

// 测试每次都分配内存的cpu、mem情况
func pprofWhenPerAlloc() { return }

// committzen， semantic-release
