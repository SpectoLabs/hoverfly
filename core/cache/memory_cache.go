package cache

import (
	"fmt"
	"sync"
)

// Cache used for storing requests and responses in memory
type InMemoryCache struct {
	elements map[string][]byte
	sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	var c InMemoryCache
	c.elements = make(map[string][]byte)
	return &c
}

func (c *InMemoryCache) Set(key, value []byte) (err error) {
	c.Lock()
	c.elements[string(key)] = value
	c.Unlock()
	return
}

func (c *InMemoryCache) Get(key []byte) (value []byte, err error) {
	c.RLock()
	bytes := c.elements[string(key)]
	value = make([]byte, len(bytes), len(bytes))
	copy(value, bytes)
	c.RUnlock()
	if len(value) == 0 {
		return nil, fmt.Errorf("key %q not found \n", key)
	}
	return value, nil
}

func (c *InMemoryCache) GetAllValues() (values [][]byte, err error) {
	c.RLock()
	values = make([][]byte, len(c.elements), len(c.elements))
	index := 0
	for _, v := range c.elements {
		values[index] = []byte(v)
		index++
	}
	c.RUnlock()
	return
}

func (c *InMemoryCache) GetAllEntries() (map[string][]byte, error) {
	c.RLock()
	dest := make(map[string][]byte)
	for k, v := range c.elements {
		dest[k] = v
	}
	c.RUnlock()
	return dest, nil
}

func (c *InMemoryCache) RecordsCount() (count int, err error) {
	c.RLock()
	len := len(c.elements)
	c.RUnlock()
	return len, nil
}

func (c *InMemoryCache) DeleteData() (err error) {
	c.Lock()
	c.elements = make(map[string][]byte)
	c.Unlock()
	return
}

func (c *InMemoryCache) GetAllKeys() (keys map[string]bool, err error) {
	c.RLock()
	keys = make(map[string]bool)
	for k, _ := range c.elements {
		keys[k] = true
	}
	c.RUnlock()
	return
}

func (c *InMemoryCache) Delete(key []byte) error {
	c.Lock()
	delete(c.elements, string(key))
	c.Unlock()
	return nil
}
