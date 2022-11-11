package mycache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_mycache/"

type HTTPPool struct {
	selfAddr string
	basePath string
}

func NewHTTPPool(selfAddr string) *HTTPPool {
	return &HTTPPool{
		selfAddr: selfAddr,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.selfAddr, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 不是指定前缀
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic(fmt.Sprintf("HTTPPool serving unexpected path: %v", r.URL.Path))
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// 约定路径为 /basePath/groupName/key
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	// 获取group
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}
	// 获取ByteView
	byteView, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(byteView.ByteSlice())
}
