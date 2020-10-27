package fast

import (
	"container/list"
	"sync"
)

type entry struct {
	key   string
	value interface{}
}
type cacheShard struct {
	locker     sync.RWMutex
	maxEntries int
	onEvicted  func(key string, value interface{})
	ll         *list.List
	cache      map[string]*list.Element
}

func newCacheShard(maxEntries int, onEvicted func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (c *cacheShard) set(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}

	en := &entry{key, value}
	e := c.ll.PushBack(en)
	c.cache[key] = e

	if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
		c.removeElement(c.ll.Front())
	}
}

func (f *cacheShard) get(key string) interface{} {
	f.locker.RLock()
	defer f.locker.RUnlock()
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

func (f *cacheShard) del(key string) {
	f.locker.RLock()
	defer f.locker.RUnlock()
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *cacheShard) delOldest() {
	f.locker.RLock()
	defer f.locker.RUnlock()
	f.removeElement(f.ll.Front())
}

func (c *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	c.ll.Remove(e)
	en := e.Value.(*entry)
	delete(c.cache, en.key)

	if c.onEvicted != nil {
		c.onEvicted(en.key, en.value)
	}
}

func (f *cacheShard) len() int {
	f.locker.RLock()
	defer f.locker.RUnlock()
	return f.ll.Len()
}
