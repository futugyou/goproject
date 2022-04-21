package code0146

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common/LinkedHashMap"
)

func Exection() {
	cache := LinkedHashMap.NewLRUCache(5)
	exection(cache)
}

func exection(cache *LinkedHashMap.LRUCache) {
	cache.Put(1, 1)
	cache.Put(2, 2)
	cache.Put(3, 3)
	cache.Put(4, 4)
	cache.Put(5, 5)
	cache.Put(6, 6)
	fmt.Println(cache.Get(2))
}
