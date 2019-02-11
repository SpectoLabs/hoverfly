package cache

import (
	"github.com/hashicorp/golang-lru"
)

// Fixed size LRU cache for storing any data type
type LRUFastCache struct {
	cache *lru.Cache
}

func NewDefaultLRUCache() *LRUFastCache {
	c, _ := lru.New(1000)
	return &LRUFastCache{ cache: c }
}

func NewLRUCache(size int) (*LRUFastCache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &LRUFastCache{ cache: c }, nil
}

func (c *LRUFastCache) Set(key, value interface{}) (err error) {
	c.cache.Add(key, value)
	return nil
}

func (c *LRUFastCache) Get(key interface{}) (value interface{}, found bool) {
	return c.cache.Get(key)
}

func (c *LRUFastCache) GetAllEntries() (map[interface{}]interface{}, error) {
	entries := make(map[interface{}]interface{}, c.cache.Len())

	for _, key := range c.cache.Keys() {
		value, _ := c.cache.Get(key)
		entries[key] = value
	}

	return entries, nil
}

func (c *LRUFastCache) RecordsCount() (count int, err error) {
	return c.cache.Len(), nil
}

func (c *LRUFastCache) DeleteData() (err error) {
	c.cache.Purge()
	return nil
}

