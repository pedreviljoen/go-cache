package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache represents a redis cache adapter implementation.
type RedisCache struct {
	c      *redis.Client
	window time.Duration
}

type cleaner struct {
	Interval time.Duration
	stop     chan bool
}

type Option func(*RedisCache)

// New -
// Initialises a new Redis client with a set of default options and passed address
func New(address, username, password string, opts ...Option) *RedisCache {
	if address == "" {
		address = "localhost:6379"
	}
	rdc := redis.NewClient(&redis.Options{
		Addr:         address,
		Username:     username,
		Password:     password,
		ReadTimeout:  time.Second * 10, // 10 second default read timeout
		WriteTimeout: time.Second * 10, // 10 second default write timeout
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Printf("redis connected")
			return nil
		},
	})
	redis := &RedisCache{
		c: rdc,
	}
	for _, opt := range opts {
		opt(redis)
	}
	return redis
}

// ClientWithCustomOptions -
// Initialises a new Redis client with provided Options
func ClientWithCustomOptions(clientOpts *redis.Options) Option {
	client := redis.NewClient(clientOpts)
	return func(rc *RedisCache) {
		rc.c = client
	}
}

// Window -
// Functional option to specify the time window of the cache
func Window(t time.Duration) Option {
	return func(rc *RedisCache) {
		rc.window = t
	}
}
