package mycache

import (
	"fmt"
	"io"
	"log"
	"mycache/consistenthash"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
)

const defaultBasePath = "/_mycache/"
const defaultReplicas = 50

type HTTPPool struct {
	selfAddr    string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map    // 一致性hash算法
	httpGetters map[string]*httpGetter // 远程节点与httpGetter映射
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

// Set 实例化一致性哈希算法, 添加节点
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 实例化一致性hash
	p.peers = consistenthash.New(defaultReplicas, nil)
	// 添加节点
	p.peers.Add(peers...)
	// 记录每个节点的httpGetter
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	// 为每个节点创建一个http客户端
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURl: peer + p.basePath}
	}
}

// PickPeer 根据key选择节点, 返回节点对应的http客户端
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 获取真实节点地址
	peer := p.peers.Get(key)
	if peer != "" && peer != p.selfAddr {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

type httpGetter struct {
	baseURl string // 要访问的远程节点地址
}

// Get 实现PeerGetter接口
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// 构建请求的url
	url := fmt.Sprintf("%v%v/%v", h.baseURl, url2.QueryEscape(group), url2.QueryEscape(key))
	// 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}
	// 返回值转byte
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}

// 确保httpGetter实现了PeerGetter接口 (类似于java的继承重写方法@Overwrite关键字的检查机制, 如果没实现下面的赋值会报错, 确实实现了但是不写下面的也没问题)
var _ PeerGetter = (*httpGetter)(nil) // 把*httpGetter类型赋值给PeerGetter, 如果没报错, 说明*httpGetter确实是实现了PeerGetter
