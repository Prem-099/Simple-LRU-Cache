package lru

type Node[K comparable, V any] struct {
	key   K
	value V
	exp   int64
	prev  *Node[K, V]
	next  *Node[K, V]
}
