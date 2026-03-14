package lru

type Node[K comparable, V any] struct {
	exp   int64
	key   K
	value V
	prev  *Node[K, V]
	next  *Node[K, V]
}
