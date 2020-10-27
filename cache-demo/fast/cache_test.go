package fast_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/goproject/cache-demo/fast"
)

func BenchmarkTourFastSetPatallel(b *testing.B) {
	cache := fast.NewFastCahe(b.N, 100000, nil)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter = counter + 1
		}
	})
}
func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}
