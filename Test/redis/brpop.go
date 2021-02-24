package main

import (
	"github.com/go-redis/redis"
	"time"
)

func main() {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0, // use default DB
	})
	res, _ := c.Ping().Result()
	println(res)

	println(time.Now().Unix())
	_, _ = c.BRPop(1e9, "l4").Result()
	//for _, v := range res1 {
	//	println(v)
	//}
	println(time.Now().UnixNano()/1e9)

}
