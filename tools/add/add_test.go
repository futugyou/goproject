package add_test

import (
	"testing"

	"github.com/goproject/tools/add"
)

func TestAdd(t *testing.T) {
	_ = add.Add("go project demo")
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add.Add("go project demo")
	}
}
