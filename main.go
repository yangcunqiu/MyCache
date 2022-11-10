package main

import (
	"fmt"
	"log"
)

func main() {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	group := New("user", 1024, GetterFunc(func(key string) ([]byte, error) {
		log.Printf("[db] search %v\n", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist\n", key)
	}))

	tom1, err := group.Get("Tom")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("key=%v, value=%v\n", "Tom", tom1)

	tom2, err := group.Get("Tom")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("key=%v, value=%v\n", "Tom", tom2)

	sam, err := group.Get("Sam")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("key=%v, value=%v\n", "Sam", sam)
}
