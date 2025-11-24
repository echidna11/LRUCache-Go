package main

import "fmt"

type Node struct {
	val  int
	next *Node
}

func (N *Node) Add(val int) bool {
	if N == nil {
		N = &Node{
			val:  val,
			next: nil,
		}
	} else {
		N.next = &Node{
			val:  val,
			next: nil,
		}
		N = N.next
	}

	return true
}


func (N *Node) Log() string {
	return fmt.Sprintf("I contain value : %d", N.val)
}

func main() {
	// LRUCache := cache.NewCache[string, int](3)
	// ok := LRUCache.Put("Vikas", 1)
	// if ok {
	// 	fmt.Println("Inserted value already exists!")
	// } else {
	// 	fmt.Println("Inserted the new entry")
	// }

	var Start Node
	fmt.Println(Start)
	Start.Add(1)
	fmt.Println(Start.Log())
}