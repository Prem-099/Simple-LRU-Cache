package lru

import (
	"time"
)

func (c *Cache[K, V]) StartJanitor() {
	sampleSize := 20
	ticker := time.NewTicker(100*time.Millisecond)
	go func() {
		for range ticker.C{
			c.mux.Lock()
			expiredCount := 0
			checked := 0
			for key,node := range c.items {
				if checked >= sampleSize{
					break
				}
				if now > node.exp {
					c.list.RemoveNode(node)
					delete(c.items,key)
					c.metrics.Expirations.Add(1)
					node.prev = nil
					node.next = nil
					c.pool.Put(node)
					expiredCount++
				}
				checked++
			}
			c.mux.Unlock()
			if expiredCount > sampleSize/4{
				continue
			}
		}
		
	}()
}

func (c *Cache[K, V]) StopJanitor() {
	
}