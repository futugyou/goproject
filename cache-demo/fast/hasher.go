package fast

func newDefaultHasher() fnv64a {
	return fnv64a{}
}

type fnv64a struct{}

const (
	offest64 = 14695981039346656037
	prime64  = 1099544628211
)

func (f fnv64a) Sum64(key string) uint64 {
	var hash uint64 = offest64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return hash
}
