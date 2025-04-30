package memory

import (
	"net/http"

	"github.com/futugyou/yomawari/semantic-kernel/abstractions/memory"
)

type MemoryBuilder struct {
	memoryStoreFactory         func() memory.IMemoryStore
	embeddingGenerationFactory func() ITextEmbeddingGenerationService
	httpClient                 *http.Client
}

func (b *MemoryBuilder) Build() *SemanticTextMemory {
	return NewSemanticTextMemory(
		b.memoryStoreFactory(),
		b.embeddingGenerationFactory(),
	)
}

type MemoryOption func(*MemoryBuilder)

func NewMemoryBuilder(opts ...MemoryOption) *MemoryBuilder {
	builder := &MemoryBuilder{}
	for _, opt := range opts {
		opt(builder)
	}
	return builder
}

func WithHttpClient(client *http.Client) MemoryOption {
	return func(b *MemoryBuilder) {
		if client == nil {
			panic("httpClient cannot be nil")
		}
		b.httpClient = client
	}
}

func WithMemoryStore(store memory.IMemoryStore) MemoryOption {
	return func(b *MemoryBuilder) {
		if store == nil {
			panic("store cannot be nil")
		}
		b.memoryStoreFactory = func() memory.IMemoryStore {
			return store
		}
	}
}

func WithMemoryStoreFactory(factory func() memory.IMemoryStore) MemoryOption {
	return func(b *MemoryBuilder) {
		if factory == nil {
			panic("factory cannot be nil")
		}
		b.memoryStoreFactory = factory
	}
}

func WithMemoryStoreFactoryWithHttp(factory func(*http.Client) memory.IMemoryStore) MemoryOption {
	return func(b *MemoryBuilder) {
		if factory == nil {
			panic("factory cannot be nil")
		}
		b.memoryStoreFactory = func() memory.IMemoryStore {
			return factory(b.httpClient)
		}
	}
}

func WithTextEmbeddingGenerationService(service ITextEmbeddingGenerationService) MemoryOption {
	return func(b *MemoryBuilder) {
		if service == nil {
			panic("store cannot be nil")
		}
		b.embeddingGenerationFactory = func() ITextEmbeddingGenerationService {
			return service
		}
	}
}

func WithTextEmbeddingGenerationServiceWithHttp(factory func(*http.Client) ITextEmbeddingGenerationService) MemoryOption {
	return func(b *MemoryBuilder) {
		if factory == nil {
			panic("factory cannot be nil")
		}
		b.embeddingGenerationFactory = func() ITextEmbeddingGenerationService {
			return factory(b.httpClient)
		}
	}
}
