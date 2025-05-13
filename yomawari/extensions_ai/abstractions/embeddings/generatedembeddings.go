package embeddings

import "github.com/futugyou/yomawari/extensions_ai/abstractions"

type GeneratedEmbeddings[TEmbedding IEmbedding] struct {
	embeddings           []TEmbedding
	Usage                *abstractions.UsageDetails
	AdditionalProperties map[string]interface{}
}

func NewGeneratedEmbeddings[TEmbedding IEmbedding]() *GeneratedEmbeddings[TEmbedding] {
	return &GeneratedEmbeddings[TEmbedding]{
		embeddings:           make([]TEmbedding, 0),
		Usage:                &abstractions.UsageDetails{},
		AdditionalProperties: map[string]interface{}{},
	}
}

func NewGeneratedEmbeddingsWithCapacity[TEmbedding IEmbedding](capacity int) *GeneratedEmbeddings[TEmbedding] {
	if capacity < 0 {
		panic("capacity cannot be less than 0")
	}
	return &GeneratedEmbeddings[TEmbedding]{
		embeddings:           make([]TEmbedding, 0, capacity),
		Usage:                &abstractions.UsageDetails{},
		AdditionalProperties: map[string]interface{}{},
	}
}

func NewGeneratedEmbeddingsFromCollection[TEmbedding IEmbedding](embeddings []TEmbedding) *GeneratedEmbeddings[TEmbedding] {
	return &GeneratedEmbeddings[TEmbedding]{
		embeddings:           append([]TEmbedding(nil), embeddings...),
		Usage:                &abstractions.UsageDetails{},
		AdditionalProperties: map[string]interface{}{},
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
		if embedding.Hash() == item.Hash() {
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
		if embedding.Hash() == item.Hash() {
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
		if embedding.Hash() == item.Hash() {
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
