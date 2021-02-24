// main退出，其他也没用了
// 非main的协程创建子协程B后退出,B是不退出的

package main

import (
	"fmt"
	"time"
)

var ch chan int

func test() int {
	ch = make(chan int)
	go func() { //父协程
		defer func() {
			fmt.Println("1")
		}()
		go func() { //子协程
			fmt.Println(<-ch)
			fmt.Println("hello")
			fmt.Println("aaaa")
		}()
	}()
	//父协程退出时，子协程没有会被强制退出啊？,还是我理解错了？
	return 0
}
func main() {
	c := test()
	time.Sleep(time.Second)
	ch <- 10
	fmt.Println("c", c)
	time.Sleep(time.Second)
}