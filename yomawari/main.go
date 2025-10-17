package main

import (
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
	"github.com/futugyou/yomawari/extensions_ai/ollama"
	"github.com/futugyou/yomawari/extensions_ai/openai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/prompts"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/search"
	"github.com/futugyou/yomawari/kernel_memory/extensions/tiktoken"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"

	coreai "github.com/futugyou/yomawari/kernel_memory/core/ai"
	coredataformats "github.com/futugyou/yomawari/kernel_memory/core/dataformats"
	coredocumentstorage "github.com/futugyou/yomawari/kernel_memory/core/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/core/filesystem"
	"github.com/futugyou/yomawari/kernel_memory/core/handlers"
	corememorystorage "github.com/futugyou/yomawari/kernel_memory/core/memorystorage"
	corepipeline "github.com/futugyou/yomawari/kernel_memory/core/pipeline"
	coreprompts "github.com/futugyou/yomawari/kernel_memory/core/prompts"
	coresearch "github.com/futugyou/yomawari/kernel_memory/core/search"
)

func main() {
	var _ embeddings.IEmbedding = (*embeddings.EmbeddingT[float64])(nil)

	var _ chatcompletion.IChatClient = (*ollama.OllamaChatClient)(nil)
	var _ embeddings.IEmbeddingGenerator[string, embeddings.EmbeddingT[float64]] = (*ollama.OllamaEmbeddingGenerator[string, embeddings.EmbeddingT[float64]])(nil)

	var _ embeddings.IEmbeddingGenerator[string, embeddings.EmbeddingT[float64]] = (*openai.OpenAIEmbeddingGenerator[embeddings.EmbeddingT[float64]])(nil)
	var _ chatcompletion.IChatClient = (*openai.OpenAIChatClient)(nil)
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
	var _ memorystorage.IMemoryDb = (*corememorystorage.SimpleVectorDb)(nil)

	var _ pipeline.IPipelineOrchestrator = (*corepipeline.BaseOrchestrator)(nil)
	var _ pipeline.IPipelineOrchestrator = (*corepipeline.DistributedPipelineOrchestrator)(nil)
	var _ pipeline.IPipelineOrchestrator = (*corepipeline.InProcessPipelineOrchestrator)(nil)

	var _ pipeline.IQueue = (*corepipeline.SimpleQueues)(nil)

	var _ search.ISearchClient = (*coresearch.SearchClient)(nil)
}
