package lru

type List[K comparable, V any] struct {
	head *Node[K, V]
	tail *Node[K, V]
}

func (l *List[K, V]) AddToFront(node *Node[K, V]) {
	node.prev = nil
	node.next = l.head
	if l.head != nil {
		l.head.prev = node
	}
	l.head = node
	if l.tail == nil {
		l.tail = node
	}
}

func (l *List[K, V]) RemoveNode(node *Node[K, V]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}
	if node.next == nil {
		l.tail = node.prev
	} else {
		node.next.prev = node.prev
	}
	node.prev = nil
	node.next = nil
}

func (l *List[K, V]) MoveToFront(node *Node[K, V]) {
	if node == l.head {
		return
	}
	l.RemoveNode(node)
	l.AddToFront(node)
}

func (l *List[K, V]) RemoveTail() *Node[K, V] {
	if l.tail == nil {
		return nil
	}
	lru := l.tail
	l.RemoveNode(lru)
	return lru
}
