package embeddings

import "github.com/futugyou/ai-extension/abstractions"

type GeneratedEmbeddings[TEmbedding comparable] struct {
	embeddings           []TEmbedding
	Usage                *abstractions.UsageDetails
	AdditionalProperties map[string]interface{}
}

func NewGeneratedEmbeddings[TEmbedding comparable]() *GeneratedEmbeddings[TEmbedding] {
	return &GeneratedEmbeddings[TEmbedding]{
		embeddings: make([]TEmbedding, 0),
	}
}

func NewGeneratedEmbeddingsWithCapacity[TEmbedding comparable](capacity int) *GeneratedEmbeddings[TEmbedding] {
	if capacity < 0 {
		panic("capacity cannot be less than 0")
	}
	return &GeneratedEmbeddings[TEmbedding]{
		embeddings: make([]TEmbedding, 0, capacity),
	}
}

func NewGeneratedEmbeddingsFromCollection[TEmbedding comparable](embeddings []TEmbedding) *GeneratedEmbeddings[TEmbedding] {
	return &GeneratedEmbeddings[TEmbedding]{
		embeddings: append([]TEmbedding(nil), embeddings...),
	}
}

func (ge *GeneratedEmbeddings[TEmbedding]) Count() int {
	return len(ge.embeddings)
}

func (ge *GeneratedEmbeddings[TEmbedding]) Add(item TEmbedding) {
	ge.embeddings = append(ge.embeddings, item)
}

func (ge *GeneratedEmbeddings[TEmbedding]) AddRange(items []TEmbedding) {
	ge.embeddings = append(ge.embeddings, items...)
}

func (ge *GeneratedEmbeddings[TEmbedding]) Clear() {
	ge.embeddings = nil
}

func (ge *GeneratedEmbeddings[TEmbedding]) Contains(item TEmbedding) bool {
	for _, embedding := range ge.embeddings {
		if embedding == item {
			return true
		}
	}
	return false
}

func (ge *GeneratedEmbeddings[TEmbedding]) Get(index int) TEmbedding {
	if index < 0 || index >= len(ge.embeddings) {
		panic("index out of bounds")
	}
	return ge.embeddings[index]
}

func (ge *GeneratedEmbeddings[TEmbedding]) Set(index int, item TEmbedding) {
	if index < 0 || index >= len(ge.embeddings) {
		panic("index out of bounds")
	}
	ge.embeddings[index] = item
}

func (ge *GeneratedEmbeddings[TEmbedding]) Remove(item TEmbedding) bool {
	for i, embedding := range ge.embeddings {
		if embedding == item {
			ge.embeddings = append(ge.embeddings[:i], ge.embeddings[i+1:]...)
			return true
		}
	}
	return false
}

func (ge *GeneratedEmbeddings[TEmbedding]) RemoveAt(index int) {
	if index < 0 || index >= len(ge.embeddings) {
		panic("index out of bounds")
	}
	ge.embeddings = append(ge.embeddings[:index], ge.embeddings[index+1:]...)
}

func (ge *GeneratedEmbeddings[TEmbedding]) IndexOf(item TEmbedding) int {
	for i, embedding := range ge.embeddings {
		if embedding == item {
			return i
		}
	}
	return -1
}

func (ge *GeneratedEmbeddings[TEmbedding]) GetEnumerator() <-chan TEmbedding {
	ch := make(chan TEmbedding)
	go func() {
		defer close(ch)
		for _, embedding := range ge.embeddings {
			ch <- embedding
		}
	}()
	return ch
}
