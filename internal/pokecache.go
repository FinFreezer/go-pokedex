package pokecache

import
(
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu *sync.Mutex
}

func NewCache(interval time.Duration) (nc *Cache){
	cache := &Cache{
		entries: make(map[string]cacheEntry),
		mu: &sync.Mutex{},
	}
	
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
	if _, ok := c.entries[key]; !ok {
		cache := cacheEntry{
			createdAt: time.Now(),
				  val: val,
		}
		c.entries[key] = cache
	}

	return
}

func (c *Cache) Get(key string) (entry []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.entries[key]; ok {
		return c.entries[key].val, true
	}

	return []byte{}, false
}

func (c *Cache) reapLoop(inter time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, val := range c.entries {
		if time.Since(val.createdAt) > inter {
			delete(c.entries, key)
		}
	}
	return
}