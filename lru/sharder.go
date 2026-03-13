package lru

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

type ShardedCache[K comparable, V any] struct {
	shards []*Cache[K, V]
	mask   uint64
}

func NewSharded[K comparable, V any](capacity int, shardCount int) *ShardedCache[K,V] {
	shards := make([]*Cache[K, V], shardCount)
	perShard := capacity / shardCount // capacity of each shard
	for i := 0; i < shardCount; i++ {
		shards[i] = New[K, V](perShard)
	}
	return &ShardedCache[K, V]{
		shards: shards,
		mask:   uint64(shardCount - 1), // shardCount should be power of 2 for better performance and efficiency
	}
}

func (s *ShardedCache[K, V]) getShard(key K) *Cache[K, V] {
	hash := fnv.New64a()
	switch v := any(key).(type) {
	case string:
		hash.Write([]byte(v))
	case int:
		b := make([]byte,8)
		binary.LittleEndian.PutUint64(b,uint64(v))
		hash.Write(b)
	default:
		hash.Write([]byte(fmt.Sprintf("%v",v)))
	}
	index := hash.Sum64() & s.mask
	return s.shards[index]
}

func (s *ShardedCache[K, V]) Get(key K) (V,bool) {
	shard := s.getShard(key)
	return shard.Get(key)
}

func (s *ShardedCache[K, V]) Put(key K,value V) {
	shard := s.getShard(key)
	shard.Put(key,value)
}