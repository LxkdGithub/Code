package main

type list struct {
	val interface{}
	first *node
	last *node
}

type node struct {
	prev *node
	next *node
	val interface{}
}

func main() {
	var a list
	a.first = nil
	//a.first = &node{
	//	val: 1,
	//}
	println(a.first.val.(int))

}
