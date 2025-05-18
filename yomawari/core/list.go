package core

type ListItem interface {
	Hash() string
}

type List[K ListItem] struct {
	data []K
}

func NewList[K ListItem]() *List[K] {
	return &List[K]{
		data: []K{},
	}
}

func (ge *List[K]) Count() int {
	return len(ge.data)
}

func (ge *List[K]) Add(item K) {
	ge.data = append(ge.data, item)
}

func (ge *List[K]) AddRange(items []K) {
	ge.data = append(ge.data, items...)
}

func (ge *List[K]) Clear() {
	ge.data = []K{}
}

func (ge *List[K]) Contains(item K) bool {
	for _, d := range ge.data {
		if d.Hash() == item.Hash() {
			return true
		}
	}
	return false
}

func (ge *List[K]) Get(index int) K {
	if index < 0 || index >= len(ge.data) {
		panic("index out of bounds")
	}
	return ge.data[index]
}

func (ge *List[K]) Items() []K {
	return ge.data
}

func (ge *List[K]) Set(index int, item K) {
	if index < 0 || index >= len(ge.data) {
		panic("index out of bounds")
	}
	ge.data[index] = item
}

func (ge *List[K]) Insert(index int, item K) {
	if index < 0 || index > len(ge.data) {
		panic("index out of bounds")
	}
	ge.data = append(ge.data, *new(K))
	copy(ge.data[index+1:], ge.data[index:])
	ge.data[index] = item
}

func (ge *List[K]) Remove(item K) bool {
	for i, d := range ge.data {
		if d.Hash() == item.Hash() {
			ge.data = append(ge.data[:i], ge.data[i+1:]...)
			return true
		}
	}
	return false
}

func (ge *List[K]) RemoveAt(index int) {
	if index < 0 || index >= len(ge.data) {
		panic("index out of bounds")
	}
	ge.data = append(ge.data[:index], ge.data[index+1:]...)
}

func (ge *List[K]) RemoveRange(index int, end int) {
	if index < 0 || end > len(ge.data) || index > end {
		panic("index out of bounds")
	}
	ge.data = append(ge.data[:index], ge.data[end:]...)
}

func (ge *List[K]) ExtractRange(index int, end int) []K {
	if index < 0 || end > len(ge.data) || index > end {
		panic("index out of bounds")
	}
	removed := make([]K, end-index)
	copy(removed, ge.data[index:end])
	ge.data = append(ge.data[:index], ge.data[end:]...)
	return removed
}

func (ge *List[K]) IndexOf(item K) int {
	for i, d := range ge.data {
		if d.Hash() == item.Hash() {
			return i
		}
	}
	return -1
}

func (ge *List[K]) GetEnumerator() <-chan K {
	ch := make(chan K)
	go func() {
		defer close(ch)
		for _, embedding := range ge.data {
			ch <- embedding
		}
	}()
	return ch
}
