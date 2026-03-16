package lru

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cache[K comparable, V any] struct {
	mux         sync.Mutex
	_           [56]byte
	items       map[K]*Node[K, V]
	list        *List[K, V]
	capacity    int
	pool        sync.Pool
	moveCounter uint64
	metrics     Metrics
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*Node[K, V], capacity),
		list:     &List[K, V]{},
		pool: sync.Pool{
			New: func() any {
				return new(Node[K, V])
			},
		},
	}
}

var now int64

func init() {
	atomic.StoreInt64(&now, time.Now().UnixNano())
	go func() {
		for {
			atomic.StoreInt64(&now, time.Now().UnixNano())
			time.Sleep(time.Millisecond)
		}
	}()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mux.Lock()
	node, ok := c.items[key]
	if !ok {
		c.mux.Unlock()
		c.metrics.Misses.Add(1)
		var zero V
		return zero, false
	}
	if node.exp > 0 && now > node.exp {
		c.list.RemoveNode(node)
		delete(c.items, key)
		node.prev = nil
		node.next = nil
		c.pool.Put(node)
		c.mux.Unlock()
		c.metrics.Expirations.Add(1)
		var zero V
		return zero, false
	}
	c.moveCounter++
	val := node.value
	c.mux.Unlock()
	c.metrics.Hits.Add(1)
	return val, true
}

func (c *Cache[K, V]) Put(key K, Value V, ttl time.Duration) {
	c.metrics.Puts.Add(1)
	exp := int64(0)
	c.mux.Lock()
	if ttl > 0 {
		exp = now + int64(ttl)
	}
	if node, ok := c.items[key]; ok {
		node.value = Value
		node.exp = exp
		c.list.MoveToFront(node)
		c.mux.Unlock()
		return
	}
	newnode := c.pool.Get().(*Node[K, V])
	newnode.key = key
	newnode.value = Value
	newnode.exp = exp
	c.list.AddToFront(newnode)
	c.items[key] = newnode
	if len(c.items) > c.capacity {
		lru := c.list.RemoveTail()
		if lru != nil {
			delete(c.items, lru.key)
			lru.prev = nil
			lru.next = nil
			c.pool.Put(lru)
			c.mux.Unlock()
			c.metrics.Evictions.Add(1)
			return
		}
	}
	c.mux.Unlock()
}

func (c *Cache[K, V]) Delete(key K) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	node, ok := c.items[key]
	if !ok {
		return false
	}
	c.list.RemoveNode(node)
	delete(c.items, node.key)
	node.next = nil
	node.prev = nil
	c.pool.Put(node)
	return true
}

func (c *Cache[K, V]) Len() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return len(c.items)
}

func (c *Cache[K, V]) Clear() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.items = make(map[K]*Node[K, V])
	c.list = &List[K, V]{}
}
