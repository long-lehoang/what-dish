package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type MemoryCache struct {
	c *gocache.Cache
}

func New(defaultExpiration, cleanupInterval time.Duration) *MemoryCache {
	return &MemoryCache{
		c: gocache.New(defaultExpiration, cleanupInterval),
	}
}

func (m *MemoryCache) Get(key string) (any, bool) {
	return m.c.Get(key)
}

func (m *MemoryCache) Set(key string, value any, ttl time.Duration) {
	m.c.Set(key, value, ttl)
}

func (m *MemoryCache) Delete(key string) {
	m.c.Delete(key)
}

func (m *MemoryCache) Flush() {
	m.c.Flush()
}
