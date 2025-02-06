package main

import (
	"fmt"
	"sync"
)

// Node -> (key : value)
type Node struct {
	key   string
	value []byte
	prev  *Node
	next  *Node
}

type LRUCache struct {
	capacity int
	cacheMap map[string]*Node
	mu       sync.Mutex
	head     *Node
	tail     *Node
}

// Returns a new node with (k,v)
func NewNode(k string, v []byte) *Node {
	return &Node{
		key:   k,
		value: v,
		next:  nil,
		prev:  nil,
	}
}

// NewLRUCache initializes with capacity
func NewLRUCache(capacity int) *LRUCache {
	lru := &LRUCache{
		capacity: capacity,
		cacheMap: make(map[string]*Node),
	}
	lru.head = NewNode("", nil)
	lru.tail = NewNode("", nil)
	lru.head.next = lru.tail
	lru.tail.prev = lru.head

	return lru
}

// Move the accessed node to the front (most recently used position)
func (lru *LRUCache) Get(key string) ([]byte, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if node, ok := lru.cacheMap[key]; ok {
		lru.removeNode(node)
		lru.addNode(node)
		return node.value, true
	}
	return nil, false
}

// Put (key,value) pair in cache
func (lru *LRUCache) Put(key string, value []byte) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	// Check if the key already exists in the cache
	if node, ok := lru.cacheMap[key]; ok {
		lru.removeNode(node)
	}

	// Create a new node and add it to the cache
	node := NewNode(key, value)
	lru.cacheMap[key] = node
	lru.addNode(node)

	// If the cache exceeds its capacity, evict the least recently used item
	if len(lru.cacheMap) > lru.capacity {
		nodeToDelete := lru.tail.prev
		lru.removeNode(nodeToDelete)
		delete(lru.cacheMap, nodeToDelete.key)
	}
}

// Add node,right after the head (most recent used position)
func (lru *LRUCache) addNode(node *Node) {
	nextNode := lru.head.next
	lru.head.next = node
	node.prev = lru.head
	node.next = nextNode
	nextNode.prev = node
}

// Remove node,left of the tail (least used position)
func (lru *LRUCache) removeNode(node *Node) {
	prevNode := node.prev
	nextNode := node.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
}

// Print the cache contents for debugging
func (lru *LRUCache) String() string {
	result := ""
	for node := lru.head.next; node != lru.tail; node = node.next {
		result += fmt.Sprintf("(%v:%v) ", node.key, node.value)
	}
	return result
}

// func main() {
// 	cache := NewLRUCache(2)
// 	cache.Put(10, 20)
// 	cache.Put(20, 30)
// 	fmt.Println("Cache after Put(10, 20) and Put(20, 30):", cache)
//
// 	v := cache.Get(10)
// 	fmt.Println("Value for key 10:", v)
// 	fmt.Println("Cache after Get(10):", cache)
//
// 	cache.Put(30, 40)
// 	fmt.Println("Cache after Put(30, 40):", cache)
//
// 	v = cache.Get(20)
// 	fmt.Println("Value for key 20:", v)
// 	fmt.Println("Cache after Get(20):", cache)
// }
