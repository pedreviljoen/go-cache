package cache

// Cache is the interface that operates the cache data.
type Cache interface {
	// Put puts value into cache with key and expire time.
	Put(key string, val []byte) error
	// Get gets cached value by given key.
	Get(key string) ([]byte, error)
	// Delete deletes cached value by given key.
	Delete(key string) error
	// IsWarm returns true if cached value exists.
	IsWarm(key string) bool
	// Flush deletes all cached data.
	Flush() error
	// FlushStale flushes all stale cached items, older than the time window
	FlushStale() error
}
