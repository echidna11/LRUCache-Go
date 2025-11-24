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

type List[T comparable, V any] struct {
	nodeMap map[T]*Node[T, V]
	mostRecent *Node[T, V] //head
	leastRecent *Node[T, V] //tail
	capacity int
	size int
}

func (l *List[T, V]) Add(key T, val V) {
	//if the item is already present in the nodeMap, call update.
	// if the item is new
	//		if the capacity is not full, just add it to the head
	//		if the capacity is full, Add it to the head, delete the tail in sequence
	if l.capacity != l.size {
		l.mostRecent.next = NewNode[T, V](key, val)
		l.mostRecent = l.mostRecent.next
		l.nodeMap[key] = l.mostRecent
	} else {
		//l.Delete()
		l.mostRecent.next = NewNode[T, V](key, val)
		l.mostRecent = l.mostRecent.next
		l.nodeMap[key] = l.mostRecent
	}
}

func (l *List[T, V]) Update(key T, val V) {
	// find the node from nodeMap
	// if the node has a left and right neighbor then connect the two
	// if the node has a right node then set that node to the tail and set this node to the head
	l.Delete(key)
	l.Add(key, val)
}

func (l *List[T, V]) Delete(key T) {
	// if left and right exist, connect the two
	// if left exists, update head
	// if right exists, update tail
	N := l.nodeMap[key]
	if N.prev != nil && N.next != nil {
		N.prev.next = N.next
		N.next.prev = N.prev
	}else if N.prev != nil {
		N = N.prev
		N.next = nil
	}else{
		N = N.next
		N.prev = nil
	}

	delete(l.nodeMap, key)
	return
}	

type Cache[K comparable, V any] struct {
	c map[K]V
	mu sync.Mutex
	list *List[K, V]
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{
		c: make(map[K]V),
		mu: sync.Mutex{},
		list: &List[K, V]{
			nodeMap: map[K]*Node[V]{},
			capacity: entryLimit,
			size: entryLimit,
		},
	}
}

// Put adds the value to the cache, and returns a boolean to indicate whether a value already existed in the cache for that key.
// If there was previously a value, it replaces that value with this one.
// Any Put counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Put(key K, value V) bool {
	_, ok := c.c[key]
	if ok {
		//update the node
		return true
	}
	c.c[key] = value
	//add the node to the linked list
	//pop out the tail node if the capacity is full 
	return false
}

// Get returns the value assocated with the passed key, and a boolean to indicate whether a value was known or not. If not, nil is returned as the value.
// Any Get counts as a refresh in terms of LRU tracking.
func (c *Cache[K, V]) Get(key K) (*V, bool) {
	v, ok := c.c[key]
	if ok {
		//update the node
		return &v, ok
	}
	return nil, ok
}


