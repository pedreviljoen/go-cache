# Go Cache

> An opinionated caching mechanism written in Go

## Usage

Using this cache implementation is fairly straight forward. The below interface are the actions that each cache adaptor will aim to implement. All supported cache adaptors will aim to satisfy the interface. This is `proto` and `JSON` safe.

```go
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
	// RunCleaner runs a process inside a go routine to flush stale cache items outside of the time window
	RunCleaner()
}
```

As an example using the in-memory caching adaptor the below code snippet instantiates a new cache and has an example of each method in the interface.

```go
package main

import (
	"encoding/json"
	"log"
	"time"

	mc "github.com/pedreviljoen/go-cache/memory"
)

type example struct {
	value   string
	number  int
}

func main() {
	c := mc.New(mc.Window(time.Minute * 5)) // instantiates a new cache with a default flush window of 5 minutes

	val := example{
		value: "Some value",
		number: 4,
    }
	cv, _ := json.Marshal(val)              // Also supports proto.Marshall for proto messages
	err := c.Put("some key", cv)            // Put some value on the cache
	if err != nil {
		log.Printf("err: %v", err)
	}
	
	val, err = c.Get("some key")            // Retrieve some value from the cache with the given key
	if err != nil {
		log.Printf("err: %v", err)
    }
	
	err = c.Delete("some key")              // Delete some value from the cache with the given key
	if err != nil {
		log.Printf("err: %v", err)
	}
	
	ok := c.IsWarm("some key")              // Checks if the key provide has a value saved in the cache
	log.Printf("result: %v", ok)
	
	err = c.FlushStale()                    // Flushes items from the cache older than the time window
	if err != nil {
		log.Printf("err: %v", err)
	}
	
	err = c.Flush()                         // Flushes all items from the cache
	if err != nil {
		log.Printf("err: %v", err)
	}
	
	c.RunCleaner()                          // Runs a cleaner process in a isolated go-routine which clears stale cache items
}
```

## Cache adaptors

- [x] In memory
- [ ] Redis
- [ ] MemCache

## Contribute

Contributions are welcome!

1. Fork it.
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

Or open up [a issue](https://github.com/pedreviljoen/go-cache/issues).

## License

[MIT License](https://github.com/pedreviljoen/go-cache/blob/main/LICENSE)