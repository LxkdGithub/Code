package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(test())
}

func test() bool {
	c := make(chan bool)
	select {
	case <-c:
		return false
	case <-time.After(10 * time.Second):
		return true
	}
}
