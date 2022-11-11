package lru

import "container/list"

type Value interface {
	Len() int
}

// Cache LRU-Cache
type Cache struct {
	maxBytes  int64                         // 最大内存大小
	nowBytes  int64                         // 当前已使用内存大小
	list      *list.List                    // 双向链表
	cache     map[string]*list.Element      // 保存key和链表中元素地址的映射
	OnEvicted func(key string, value Value) // key被删除后的回调函数
}

type entry struct {
	key   string
	value Value
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Set 增加
func (c *Cache) Set(key string, value Value) {
	element, ok := c.cache[key]
	if ok {
		// 修改
		// key被使用, 移到到队尾 (list是双向链表, 所以我们约定Front的方向是队尾)
		c.list.MoveToFront(element)
		// 链表元素强转为entry
		entry := element.Value.(*entry)
		// 重新计算已用内存大小
		c.nowBytes += int64(value.Len()) - int64(entry.value.Len())
		// 修改value
		entry.value = value
	} else {
		// 新增
		// 新元素放在链表队尾
		element := c.list.PushFront(&entry{key: key, value: value})
		// 保存map
		c.cache[key] = element
		// 计算已用内存大小
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	// 新增后判断已用内存是否超过最大内存限制
	if c.maxBytes != 0 && c.nowBytes > c.maxBytes {
		// 移除最近最少使用 lru
		c.RemoveOldest()
	}
}

// RemoveOldest lru实现
func (c *Cache) RemoveOldest() {
	// 获取队首元素
	back := c.list.Back()
	if back != nil {
		// 移除队首元素
		c.list.Remove(back)
		// 链表元素强转成entry
		entry := back.Value.(*entry)
		// 删除map中key和链表元素的映射
		delete(c.cache, entry.key)
		// 减小已用内存大小
		c.nowBytes -= int64(len(entry.key)) + int64(entry.value.Len())
		// 执行回调
		if c.OnEvicted != nil {
			c.OnEvicted(entry.key, entry.value)
		}
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	element, ok := c.cache[key]
	if !ok {
		return
	}
	// 这个key使用了一次, 将对应元素移到链表队尾
	c.list.MoveToFront(element)
	entry := element.Value.(*entry)
	return entry.value, true
}

func (c *Cache) Len() int {
	return c.list.Len()
}
