# MyCache

## 目录结构
````
lru/
    |--lru.go  // lru 缓存淘汰策略
byteview.go // 缓存值的抽象与封装
cache.go    // 并发控制
mycache.go // 负责与外部交互，控制缓存存储和获取的主流程
http.go  // 提供被其他节点访问的能力
````