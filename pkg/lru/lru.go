package lru

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Item struct {
	Key      int32
	Value    string
	CreateAt time.Time
}

type Cache struct {
	size   int
	expire time.Duration
	mu     sync.RWMutex
	items  map[int32]*list.Element
	queue  *list.List
}

func NewCache(size int) *Cache {
	return &Cache{
		size:   size,
		expire: time.Duration(0),
		items:  make(map[int32]*list.Element),
		queue:  list.New(),
	}
}

func NewTTLCache(size int, ttl time.Duration) *Cache {
	return &Cache{
		size:   size,
		expire: ttl,
		items:  make(map[int32]*list.Element),
		queue:  list.New(),
	}
}

func (c *Cache) Put(key int32, value string) bool {
	c.mu.RLock()
	element, exists := c.items[key]
	c.mu.RUnlock()

	if exists {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		return true
	}

	if c.queue.Len() == c.size {
		if element := c.queue.Back(); element != nil {
			c.deleteElement(element)
		}
	}

	var createAt time.Time
	if c.expire != 0 {
		createAt = time.Now()
	}
	item := &Item{
		Key:      key,
		Value:    value,
		CreateAt: createAt,
	}

	element = c.queue.PushFront(item)

	c.mu.Lock()
	c.items[item.Key] = element
	c.mu.Unlock()
	return true
}

func (c *Cache) Get(key int32) string {
	c.mu.RLock()
	element, exists := c.items[key]
	c.mu.RUnlock()
	if !exists {
		return "there is no such key"
	}

	item := element.Value.(*Item)

	if c.expire != 0 && time.Now().Sub(item.CreateAt) > c.expire {
		c.deleteElement(element)
		return "not found"

	}
	c.queue.MoveToFront(element)
	return item.Value
}

func (c *Cache) deleteElement(element *list.Element) {
	c.mu.Lock()
	item := c.queue.Remove(element).(*Item)
	delete(c.items, item.Key)
	c.mu.Unlock()
}

func (c *Cache) ShowCurrentCache() {
	var res string
	res += fmt.Sprintf("{")
	for key, val := range c.items {
		res += fmt.Sprintf(`%v: "%v", `, key, val.Value.(*Item).Value)
	}
	if len(c.items) > 0 {
		res = res[0 : len(res)-2]
	}
	res += fmt.Sprintf("}")
	fmt.Println(res)
}
