package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu sync.RWMutex

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value, true
	}

	return nil, false
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		delete(c.items, key)
	}

	c.items[key] = c.queue.PushFront(value)

	return false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.capacity = 0
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
