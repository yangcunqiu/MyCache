package main

import (
	"flag"
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

func createGroup() *mycache.Group {
	return mycache.New("user", 1024, mycache.GetterFunc(func(key string) ([]byte, error) {
		log.Printf("[db] search %v\n", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist\n", key)
	}))
}

func startCacheServer(addr string, addrs []string, mc *mycache.Group) {
	peers := mycache.NewHTTPPool(addr)
	peers.Set(addrs...)
	mc.RegisterPeers(peers)
	log.Println("mycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, mc *mycache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		byteView, err := mc.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(byteView.ByteSlice())
	}))
	log.Println("fontend server is running ad:", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	mc := createGroup()
	if api {
		go startAPIServer(apiAddr, mc)
	}
	startCacheServer(addrMap[port], []string(addrs), mc)
}
