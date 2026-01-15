package pokecache

import
(
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	Val []byte
}

type Cache struct {
	Entries map[string]cacheEntry
	mu *sync.Mutex
}

func NewCache(interval time.Duration) (nc *Cache){
	cache := &Cache{
		Entries: make(map[string]cacheEntry),
		mu: &sync.Mutex{},
	}
	cache.reapLoop(interval)
	go cache.timer(interval)
	return cache
}

func (c *Cache) timer(interval time.Duration) {
	ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for range ticker.C {
        c.reapLoop(interval)
    }
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if _, ok := c.Entries[key]; !ok {
		cache := cacheEntry{
			createdAt: time.Now(),
				  Val: val,
		}
		c.Entries[key] = cache
	}

	return
}

func (c *Cache) Get(key string) (entry []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[key]; ok {
		return c.Entries[key].Val, true
	}

	return []byte{}, false
}

func (c *Cache) reapLoop(inter time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, val := range c.Entries {
		if time.Since(val.createdAt) > inter {
			delete(c.Entries, key)
		}
	}
	return
}