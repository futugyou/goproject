package cache

import (
	co "github.com/futugyousuzu/goproject/algorithm/common/LinkedHashMap"
)

type LRUCache struct {
	dic   map[int]*co.Node
	cache *co.DoubleList
	cap   int
}

func NewLRUCache(cap int) *LRUCache {
	l := LRUCache{}
	l.cap = cap
	l.cache = co.NewDoubleList()
	l.dic = make(map[int]*co.Node)
	return &l
}

func (l *LRUCache) makeRecently(key int) {
	x := l.dic[key]
	l.cache.Remove(x)
	l.cache.AddLast(x)
}

func (l *LRUCache) addRecently(key, val int) {
	x := co.NewNode(key, val)
	l.cache.AddLast(x)
	l.dic[key] = x
}

func (l *LRUCache) deleteKey(key int) {
	x := l.dic[key]
	l.cache.Remove(x)
	delete(l.dic, key)
}

func (l *LRUCache) removeLeastRecently() {
	deletedNode := l.cache.RemoveFirst()
	deletedKey := deletedNode.Key
	delete(l.dic, deletedKey)
}

func (l *LRUCache) Get(key int) int {
	if _, ok := l.dic[key]; ok {
		l.makeRecently(key)
		return l.dic[key].Value
	}
	return -1
}
func (l *LRUCache) Put(key, val int) {
	if _, ok := l.dic[key]; ok {
		l.deleteKey(key)
		l.addRecently(key, val)
		return
	}
	if l.cap == l.cache.Size() {
		l.removeLeastRecently()
	}
	l.addRecently(key, val)
}
