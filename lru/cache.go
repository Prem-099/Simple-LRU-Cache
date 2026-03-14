package lru

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*Node[K, V]
	list     *List[K, V]
	mux      sync.Mutex
	pool     sync.Pool
}

func New[K comparable,V any] (capacity int) *Cache[K,V] {
	return &Cache[K,V]{
		capacity: capacity,
		items: make(map[K]*Node[K, V],capacity),
		list: &List[K, V]{},
		pool: sync.Pool{
			New: func() any {
				return new(Node[K,V])
			},
		},

	}
}

var now int64
func init() {
	atomic.StoreInt64(&now,time.Now().UnixNano())
	go func(){
		for{
			atomic.StoreInt64(&now,time.Now().UnixNano())
			time.Sleep(time.Millisecond)
		}
	}()
}


func (c *Cache[K, V]) Get(Key K) (V, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if node,ok := c.items[Key]; ok{
		cnow := atomic.LoadInt64(&now)
		if cnow > node.exp{
			c.list.RemoveNode(node)
			delete(c.items,Key)

			var zero V
			return zero,false
		}
		c.list.MoveToFront(node)
		return node.value,true
	}
	var zero V
	return zero,false
}

func (c *Cache[K, V]) Put(Key K, Value V,ttl time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	exp := int64(0)
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}else{
		exp = time.Now().Add(5*time.Second).UnixNano()
	}
	if node,ok := c.items[Key];ok{
		node.value = Value
		node.exp = exp
		c.list.MoveToFront(node)
		return
	}
	newnode:=c.pool.Get().(*Node[K,V])
	newnode.key = Key
	newnode.value = Value
	newnode.exp = exp
	c.list.AddToFront(newnode)
	c.items[Key] = newnode
	if len(c.items) > c.capacity{
		lru := c.list.RemoveTail()
		if lru!=nil{
			delete(c.items,lru.key)
			lru.prev = nil
			lru.next = nil
			c.pool.Put(lru)
		}
	}
}