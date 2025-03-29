package main

import (
	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/generative-ai/abstractions/embeddings"
	"github.com/futugyou/yomawari/generative-ai/ollama"
	"github.com/futugyou/yomawari/generative-ai/openai"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/prompts"
	"github.com/futugyou/yomawari/kernel-memory/extensions/tiktoken"

	"github.com/futugyou/yomawari/kernel-memory/abstractions/ai"

	corememorystorage "github.com/futugyou/yomawari/kernel-memory/core/memorystorage"
	coreai "github.com/futugyou/yomawari/kernel-memory/core/ai"
	coredataformats "github.com/futugyou/yomawari/kernel-memory/core/dataformats"
	coredocumentstorage "github.com/futugyou/yomawari/kernel-memory/core/documentstorage"
	"github.com/futugyou/yomawari/kernel-memory/core/filesystem"
	"github.com/futugyou/yomawari/kernel-memory/core/handlers"
	coreprompts "github.com/futugyou/yomawari/kernel-memory/core/prompts"
)

func main() {
	var _ embeddings.IEmbedding = (*embeddings.EmbeddingT[float64])(nil)

	var _ chatcompletion.IChatClient = (*ollama.OllamaChatClient)(nil)
	var _ embeddings.IEmbeddingGenerator[string, embeddings.EmbeddingT[float64]] = (*ollama.OllamaEmbeddingGenerator[string, embeddings.EmbeddingT[float64]])(nil)

	var _ embeddings.IEmbeddingGenerator[string, embeddings.EmbeddingT[float64]] = (*openai.OpenAIEmbeddingGenerator[embeddings.EmbeddingT[float64]])(nil)
	var _ chatcompletion.IChatClient = (*openai.OpenAIChatClient)(nil)
	var _ chatcompletion.IChatClient = (*openai.OpenAIAssistantClient)(nil)
	var _ chatcompletion.IChatClient = (*openai.OpenAIResponseChatClient)(nil)

	var _ ai.ITextEmbeddingGenerator = (*coreai.NoEmbeddingGenerator)(nil)
	var _ ai.ITextGenerator = (*coreai.NoTextGenerator)(nil)
	var _ ai.ITextTokenizer = (*tiktoken.CL100KTokenizer)(nil)

	var _ context.IContext = &context.RequestContext{}
	var _ context.IContextProvider = &context.RequestContextProvider{}

	var _ dataformats.IContentDecoder = (*coredataformats.HtmlDecoder)(nil)
	var _ dataformats.IContentDecoder = (*coredataformats.ImageDecoder)(nil)
	var _ dataformats.IContentDecoder = (*coredataformats.MarkDownDecoder)(nil)
	var _ dataformats.IContentDecoder = (*coredataformats.MsExcelDecoder)(nil)
	var _ dataformats.IContentDecoder = (*coredataformats.MsPowerPointDecoder)(nil)
	var _ dataformats.IContentDecoder = (*coredataformats.MsWordDecoder)(nil)
	var _ dataformats.IContentDecoder = (*coredataformats.PdfDecoder)(nil)
	var _ dataformats.IContentDecoder = (*coredataformats.TextDecoder)(nil)
	var _ dataformats.IWebScraper = (*coredataformats.WebScraper)(nil)

	var _ documentstorage.IDocumentStorage = (*coredocumentstorage.SimpleFileStorage)(nil)

	var _ pipeline.IMimeTypeDetection = (*pipeline.MimeTypesDetection)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.DeleteDocumentHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.DeleteGeneratedFilesHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.DeleteIndexHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.GenerateEmbeddingsHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.GenerateEmbeddingsParallelHandler)(nil)

	var _ pipeline.IPipelineStepHandler = (*handlers.SaveRecordsHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.SummarizationHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.SummarizationParallelHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.TextExtractionHandler)(nil)
	var _ pipeline.IPipelineStepHandler = (*handlers.TextPartitioningHandler)(nil)

	var _ prompts.IPromptProvider = (*coreprompts.EmbeddedPromptProvider)(nil)

	var _ filesystem.IFileSystem = (*filesystem.DiskFileSystem)(nil)
	var _ filesystem.IFileSystem = (*filesystem.VolatileFileSystem)(nil)

	var _ memorystorage.IMemoryDb = (*corememorystorage.SimpleTextDb)(nil)

}
