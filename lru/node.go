package lru

type Node[K comparable, V any] struct{
	key K
	value V
	prev *Node[K, V]
	next *Node[K, V]
}