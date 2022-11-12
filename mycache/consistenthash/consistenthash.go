package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map 一致性hash主结构
type Map struct {
	hash     Hash           // hash
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 虚拟节点与真实节点的映射表 key:虚拟节点hash值, value:真实节点
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hash,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		// 使用默认hash
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点 key: 真实节点名称
func (m *Map) Add(keys ...string) {
	// 为多个节点添加
	for _, key := range keys {
		// 根据虚拟节点倍数去添加
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点名称=编号+真实节点名称
			virtualKey := strconv.Itoa(i) + key
			// 对虚拟节点hash
			hash := int(m.hash([]byte(virtualKey)))
			// 添加到环上
			m.keys = append(m.keys, hash)
			// 记录虚拟节点和真实节点的映射
			m.hashMap[hash] = key
		}
	}
	// 哈希环排序
	sort.Ints(m.keys)
}

// Get 获取节点 key:要保存的key, string:真实节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 对key进行hash计算
	hash := int(m.hash([]byte(key)))
	// 顺时针寻找环上第一个匹配的虚拟节点下标
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 因为m.keys是环状结构, 所以下标应该取余
	index := idx % len(m.keys)
	// 获取下标对应的虚拟节点hash值
	virtualHash := m.keys[index]
	// 通过映射表获取真实节点
	return m.hashMap[virtualHash]
}
