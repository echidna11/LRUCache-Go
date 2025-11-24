package cache

import "sync"

type Node[T comparable, V any] struct {
	next *Node[T, V]
	prev *Node[T, V]
	value V
	key T
}

func NewNode[T comparable, V any](key T, val V) *Node[T, V] {
	return &Node[T, V]{value: val, key: key}
} 

func (N *Node[T, V]) GetValue() *V {
	return &N.value
}

type List[T comparable, V any] struct {
	MostRecent *Node[T, V] //head
	LeastRecent *Node[T, V] //tail
	capacity int
	size int
}

func (l *List[T, V]) Insert(N *Node[T, V]) {
	if l.capacity == 0 {
		l.LeastRecent = N
		l.MostRecent = N
		l.capacity += 1
		return
	}

	l.MostRecent.next = N
	N.next = nil
	N.prev = l.MostRecent
	l.MostRecent = l.MostRecent.next
	l.capacity += 1

	if l.capacity > l.size {
		l.LeastRecent = l.LeastRecent.next
		l.LeastRecent.prev = nil
		l.capacity -= 1
	}
	return
}

func (l *List[T, V]) Update(OldNode, NewNode *Node[T, V] ) {
	l.Delete(OldNode)
	l.Insert(NewNode)
}

func (l *List[T, V]) Max() bool {
	return l.capacity == l.size
}

func (l *List[T, V]) Tail() T {
	return l.LeastRecent.key
}
 
func (l *List[T, V]) Delete(N *Node[T,V]) {
	if N.prev != nil && N.next != nil {
		N.prev.next = N.next
		N.next.prev = N.prev
	}else if N.prev != nil {
		l.MostRecent = N.prev
		N.prev.next = nil
	}else if N.next != nil {
		l.LeastRecent = l.LeastRecent.next
		l.LeastRecent.prev = nil
	}else {
		l.LeastRecent = nil
		l.MostRecent = nil
	}
	l.capacity -= 1
	return
}	

type Cache[T comparable, V any] struct {
	c map[T]*Node[T, V]
	mu sync.Mutex
	list *List[T, V]
}

func NewCache[T comparable, V any](entryLimit int) Cache[T, V] {
	return Cache[T, V]{
		c: make(map[T]*Node[T, V]),
		mu: sync.Mutex{},
		list: &List[T, V]{
			capacity: 0,
			size: entryLimit,
		},
	}
}

func (c *Cache[T, V]) Put(key T, value V) bool {
	N, ok := c.c[key]
	if ok {
		M := NewNode[T, V](key, value)
		c.c[key] = M
		c.list.Update(N, M)
		return true
	}
	N = NewNode[T, V](key, value)
	c.c[key] = N
	if c.list.Max() {
		delete(c.c, c.list.Tail())
	}
	c.list.Insert(N)
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	N, ok := c.c[key]
	if ok {
		c.list.Update(N, N)
		return N.GetValue(), ok
	}
	return nil, ok
}


