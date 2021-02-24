package main

import "sync/atomic"

type A struct {
	a int32
}

func main() {
	var A1 = A{a: 1}
	println(A1.a)
	atomic.AddInt32(&A1.a, 1)
	a := "sdsd"
	println(&(a[1]))
	b := []byte("sdsd")
	println(&b[1])
}
