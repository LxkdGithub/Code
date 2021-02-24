package main

import (
)

func main() {
	var a, b interface{}
	a = []byte("ls")
	b = []byte("lsk")
	//println(reflect.TypeOf(a) == reflect.TypeOf(b))
	//println(bytes.Compare(a.([]byte), b.([]byte)))
	a = 1
	b = 2
	println(a)
	println(b.(int))
}
