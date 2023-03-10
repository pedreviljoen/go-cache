package redis

import (
	"context"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
)

// IsWarm -
// Accept a cache key identifier and determines if the cache is still within
// the time duration window
func (c *RedisCache) IsWarm(key string) bool {
	_, err := c.c.Exists(context.Background(), key).Result()
	return err != redis.Nil
}

// Put -
// Accepts a cache key identifier and value, save the respective key and value
// inside the Redis cache
func (c *RedisCache) Put(key string, value []byte) error {
	if err := c.c.Set(context.Background(), key, value, c.window).Err(); err != nil {
		return err
	}
	return nil
}

// Get -
// Accepts a cache key identifier and fetches the value of the corresponding cache key
func (c *RedisCache) Get(key string) ([]byte, error) {
	val, err := c.c.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

// Delete -
// Accepts a cache item key identifier and deletes the value of the corresponding cache key
func (c *RedisCache) Delete(key string) error {
	if err := c.c.Del(context.Background(), key).Err(); err != nil {
		return err
	}
	return nil
}

// Flush -
// Empties the entire cache
func (c *RedisCache) Flush() error {
	ctx := context.Background()
	iter := c.c.Scan(ctx, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := c.c.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

// FlushStale -
// Iterates over all cache key-value items and removes all stale cache items
func (c *RedisCache) FlushStale() error {
	ctx := context.Background()
	iter := c.c.Scan(ctx, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		d, err := c.c.TTL(ctx, key).Result()
		if err != nil {
			return err
		}

		if d == -1 { // -1 means no TTL
			if err := c.c.Del(ctx, key).Err(); err != nil {
				return err
			}
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

// RunCleaner -
// Initialises and starts a new cleaner process in a separate go routine
// this process flushes cache items inside the cache which are older than the configured cache window
func (c *RedisCache) RunCleaner() {
	j := c.initCleaner()
	j.run(c)
	runtime.SetFinalizer(c, j.stopCleaner)
}

// initCleaner -
// Initialises a new cleaner
func (c *RedisCache) initCleaner() *cleaner {
	return &cleaner{
		Interval: c.window,
		stop:     make(chan bool),
	}
}

// runCleaner -
// Runs the cleaner inside a go routine
func (j *cleaner) run(c *RedisCache) {
	go j.cleanup(c)
}

// cleanup -
// Calls the underlying FlushStale method of the cache which clears
// stale cache items
func (j *cleaner) cleanup(c *RedisCache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.FlushStale()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

// stopCleaner -
// Sends a stop signal to the go-routine running the cleaner process
func (j *cleaner) stopCleaner() {
	j.stop <- true
}
