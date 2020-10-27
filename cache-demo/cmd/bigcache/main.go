package main

import (
	"log"
	"time"

	"github.com/allegro/bigcache/v2"
)

func main() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(20 * time.Minute))
	if err != nil {
		log.Printf("one %v", err)
		return
	}

	// if key not exsits ,there will be an Entry not found error
	entry, err := cache.Get("my-unique-key")
	if err != nil && entry != nil {
		log.Printf("two %v", err)
		return
	}

	if entry == nil {
		entry = []byte("value")
		cache.Set("my-unique-key", entry)
	}
	log.Println(string(entry))
}

type BigCache struct {
	Shards     int
	lifeWindow uint64
}

type BigCacheOption func(*BigCache)

func ShardsNum(shards int) BigCacheOption {
	return func(c *BigCache) {
		c.Shards = shards
	}
}

func LifeWindow(eviction time.Duration) BigCacheOption {
	return func(c *BigCache) {
		c.lifeWindow = uint64(eviction.Seconds())
	}
}

func NewBigCache(options ...BigCacheOption) (*BigCache, error) {
	c := &BigCache{}
	for _, f := range options {
		f(c)
	}
	return c, nil
}
