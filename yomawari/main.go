package main

import (
	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/generative-ai/abstractions/embeddings"
	"github.com/futugyou/yomawari/generative-ai/ollama"
	"github.com/futugyou/yomawari/generative-ai/openai"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/pipeline"

	"github.com/futugyou/yomawari/kernel-memory/abstractions/ai"

	coreai "github.com/futugyou/yomawari/kernel-memory/core/ai"
	"github.com/futugyou/yomawari/kernel-memory/core/filesystem"
)

func main() {
	var _ embeddings.IEmbedding = (*embeddings.EmbeddingT[float64])(nil)

	var _ chatcompletion.IChatClient = (*ollama.OllamaChatClient)(nil)
	var _ embeddings.IEmbeddingGenerator[string, embeddings.EmbeddingT[float64]] = (*ollama.OllamaEmbeddingGenerator[string, embeddings.EmbeddingT[float64]])(nil)

	var _ embeddings.IEmbeddingGenerator[string, embeddings.EmbeddingT[float64]] = (*openai.OpenAIEmbeddingGenerator[embeddings.EmbeddingT[float64]])(nil)
	var _ chatcompletion.IChatClient = (*openai.OpenAIChatClient)(nil)
	var _ chatcompletion.IChatClient = (*openai.OpenAIAssistantClient)(nil)
	var _ chatcompletion.IChatClient = (*openai.OpenAIResponseChatClient)(nil)

	var _ context.IContext = &context.RequestContext{}
	var _ context.IContextProvider = &context.RequestContextProvider{}
	var _ pipeline.IMimeTypeDetection = &pipeline.MimeTypesDetection{}
	var _ ai.ITextEmbeddingGenerator = (*coreai.NoEmbeddingGenerator)(nil)

	var _ filesystem.IFileSystem = (*filesystem.DiskFileSystem)(nil)
	var _ filesystem.IFileSystem = (*filesystem.VolatileFileSystem)(nil)
}
