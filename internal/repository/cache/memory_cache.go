package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value      interface{}
	ExpiresAt  time.Time
}

type MemoryCache struct {
	data sync.Map
	mu   sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{}
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	expiresAt := time.Now().Add(ttl)
	entry := CacheEntry{
		Value:     value,
		ExpiresAt: expiresAt,
	}
	c.data.Store(key, entry)
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	val, exists := c.data.Load(key)
	if !exists {
		return nil, false
	}

	entry := val.(CacheEntry)
	if time.Now().After(entry.ExpiresAt) {
		// Entry has expired, remove it
		c.data.Delete(key)
		return nil, false
	}

	return entry.Value, true
}

func (c *MemoryCache) Delete(key string) {
	c.data.Delete(key)
}

func (c *MemoryCache) Clear() {
	c.data.Range(func(key, value interface{}) bool {
		c.data.Delete(key)
		return true
	})
}

// CleanupExpired removes all expired entries
func (c *MemoryCache) CleanupExpired() {
	c.data.Range(func(key, value interface{}) bool {
		entry := value.(CacheEntry)
		if time.Now().After(entry.ExpiresAt) {
			c.data.Delete(key)
		}
		return true
	})
}
