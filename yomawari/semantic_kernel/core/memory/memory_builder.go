package memory

import (
	"net/http"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

type MemoryBuilder struct {
	memoryStoreFactory         func() abstractions.IMemoryStore
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

func WithMemoryStore(store abstractions.IMemoryStore) MemoryOption {
	return func(b *MemoryBuilder) {
		if store == nil {
			panic("store cannot be nil")
		}
		b.memoryStoreFactory = func() abstractions.IMemoryStore {
			return store
		}
	}
}

func WithMemoryStoreFactory(factory func() abstractions.IMemoryStore) MemoryOption {
	return func(b *MemoryBuilder) {
		if factory == nil {
			panic("factory cannot be nil")
		}
		b.memoryStoreFactory = factory
	}
}

func WithMemoryStoreFactoryWithHttp(factory func(*http.Client) abstractions.IMemoryStore) MemoryOption {
	return func(b *MemoryBuilder) {
		if factory == nil {
			panic("factory cannot be nil")
		}
		b.memoryStoreFactory = func() abstractions.IMemoryStore {
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
