package main

import (
	"MyCache/cache"
	"fmt"
	"log"
)

func main() {
	c := cache.New(20, func(key string, value cache.Value) {
		log.Printf("delete key=%v", key)
	})
	c.Set("name1", cache.String("ycq"))
	c.Set("name2", cache.String("zml"))
	c.Set("name3", cache.String("zml"))

	value1, ok1 := c.Get("name1")
	value2, ok2 := c.Get("name2")
	fmt.Printf("key=%v, ok=%v, value=%v\n", "name1", ok1, value1)
	fmt.Printf("key=%v, ok=%v, value=%v\n", "name2", ok2, value2)

}
