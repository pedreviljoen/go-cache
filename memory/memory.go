package memory

import (
	"sync"
	"time"
)

const defaultWindow = time.Second * 60

// MemCache is an in memory cache implementation
type MemCache struct {
	mutex  sync.RWMutex
	window time.Duration
	cache  map[string]MemCacheValue
}

// MemCacheValue represents a cached value as part of MemCache
type MemCacheValue struct {
	saved time.Time // when this value was saved
	value []byte    // result of proto.Marshal()
}

type Option func(*MemCache)

// New -
// Constructor function which initialises a new cache
// accepts a time duration window and cache key identifier separator
func New(opts ...Option) *MemCache {
	nache := &MemCache{
		cache:  map[string]MemCacheValue{},
		mutex:  sync.RWMutex{},
		window: defaultWindow,
	}
	for _, opt := range opts {
		opt(nache)
	}
	return nache
}

// Window -
// Functional option to specify the time window of the cache
func Window(t time.Duration) Option {
	return func(mc *MemCache) {
		mc.window = t
	}
}
