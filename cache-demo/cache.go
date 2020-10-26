package cache

import (
	"log"
	"sync"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

const DefaultMaxBytes = 1 << 29

type safeCache struct {
	m          sync.RWMutex
	cache      Cache
	nhit, nget int
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (sc *safeCache) set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, value)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.nget++
	if sc.cache == nil {
		return nil
	}

	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[Cache] hit")
		sc.nhit++
	}
	return v
}

func (sc *safeCache) stat() *Stat {
	sc.m.Lock()
	defer sc.m.Unlock()
	return &Stat{
		NHit: sc.nhit,
		NGet: sc.nget,
	}
}

type Stat struct {
	NHit, NGet int
}
