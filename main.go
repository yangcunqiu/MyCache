package main

import (
	"fmt"
	"log"
	"mycache"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	mycache.New("user", 1024, mycache.GetterFunc(func(key string) ([]byte, error) {
		log.Printf("[db] search %v\n", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist\n", key)
	}))

	addr := "localhost:9999"
	httpPool := mycache.NewHTTPPool(addr)
	log.Printf("mycache is running at %v", addr)
	log.Fatal(http.ListenAndServe(addr, httpPool))
}
