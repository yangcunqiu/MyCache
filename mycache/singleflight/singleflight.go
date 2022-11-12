package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 对于相同key, 无论Do被调用多少次, fn只会被调用一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		// 如果key对应的请求已经存在
		g.mu.Unlock()
		c.wg.Wait()         // 等待请求执行完成
		return c.val, c.err // 请求结束, 返回结果
	}
	// key对应请求不存在
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c // 表示key对应的请求已经被处理
	g.mu.Unlock()

	// 调用fn 发起请求
	c.val, c.err = fn() // 调用过程中, 相同key进来会等待
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	// 请求结束, 删掉请求
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
