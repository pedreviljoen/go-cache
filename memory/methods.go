package memory

import (
	"fmt"
	"time"
)

// IsWarm -
// Accept a cache key identifier and determines if the cache is still within
// the time duration window
func (c *MemCache) IsWarm(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok := c.cache[key]
	age := (time.Since(val.saved) - c.window) * -1
	return ok && age > 0
}

// Put -
// Accepts a cache key identifier and value, save the respective key and value
// inside the in-memory cache
func (c *MemCache) Put(key string, value []byte) error {
	cache := map[string]MemCacheValue{}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	curCache := c.cache
	for k, v := range curCache {
		if k != key {
			cache[k] = v
		}
	}
	nVal := MemCacheValue{
		value: value,
		saved: time.Now(),
	}
	cache[key] = nVal
	c.cache = cache
	return nil
}

// Get -
// Accepts a cache key identifier and fetches the value of the corresponding cache key
func (c *MemCache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	cache, ok := c.cache[key]
	if !ok {
		return nil, fmt.Errorf("unable to retrieve value from cache")
	}
	return cache.value, nil
}

// Delete -
// Accepts a cache key identifier and deletes the value of the corresponding cache key
func (c *MemCache) Delete(key string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.cache[key]
	if !ok {
		return fmt.Errorf("unable to retrieve value from cache")
	} else {
		c.cache[key] = MemCacheValue{}
	}
	return nil
}

// Flush -
// Empties the entire cache
func (c *MemCache) Flush() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	nCache := map[string]MemCacheValue{}
	c.cache = nCache
	return nil
}

// FlushStale -
// Iterates over all cache key-value items and removes all stale cache items
func (c *MemCache) FlushStale() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.cache {
		age := (time.Since(v.saved) - c.window) * (-1)
		if age < 0 {
			delete(c.cache, k)
		}
	}
	return nil
}
