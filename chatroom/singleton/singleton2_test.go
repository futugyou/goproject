package singleton

import (
	"testing"
)

//go.exe test -benchmem -run=^$ github.com/goproject/chatroom/singleton -bench ^(BenchmarkNew)$
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}
