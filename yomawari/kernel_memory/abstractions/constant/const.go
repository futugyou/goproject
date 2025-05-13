package constant

const (
	// WebService Constants
	IndexField      = "index"
	DocumentIdField = "documentId"
	FilenameField   = "filename"
	TagsField       = "tags"
	StepsField      = "steps"
	ArgsField       = "args"

	// ModelType Constants
	EmbeddingGeneration = "EmbeddingGeneration"
	TextGeneration      = "TextGeneration"

	// Pipeline status
	PipelineStatusFilename = "__pipeline_status.json"

	// Tags settings
	ReservedEqualsChar             = ':'
	ReservedTagsPrefix             = "__"
	ReservedDocumentIdTag          = ReservedTagsPrefix + "document_id"
	ReservedFileIdTag              = ReservedTagsPrefix + "file_id"
	ReservedFilePartitionTag       = ReservedTagsPrefix + "file_part"
	ReservedFilePartitionNumberTag = ReservedTagsPrefix + "part_n"
	ReservedFileSectionNumberTag   = ReservedTagsPrefix + "sect_n"
	ReservedFileTypeTag            = ReservedTagsPrefix + "file_type"
	ReservedSyntheticTypeTag       = ReservedTagsPrefix + "synth"

	// Known tags
	TagsSyntheticSummary = "summary"

	// Payload fields
	ReservedPayloadSchemaVersionField   = "schema"
	ReservedPayloadTextField            = "text"
	ReservedPayloadFileNameField        = "file"
	ReservedPayloadUrlField             = "url"
	ReservedPayloadLastUpdateField      = "last_update"
	ReservedPayloadVectorProviderField  = "vector_provider"
	ReservedPayloadVectorGeneratorField = "vector_generator"

	// Endpoints
	HttpAskEndpoint                      = "/ask"
	HttpSearchEndpoint                   = "/search"
	HttpDownloadEndpoint                 = "/download"
	HttpUploadEndpoint                   = "/upload"
	HttpUploadStatusEndpoint             = "/upload-status"
	HttpDocumentsEndpoint                = "/documents"
	HttpIndexesEndpoint                  = "/indexes"
	HttpIndexPlaceholder                 = "{index}"
	HttpDocumentIdPlaceholder            = "{documentId}"
	HttpFilenamePlaceholder              = "{filename}"
	HttpDeleteDocumentEndpointWithParams = "{HttpDocumentsEndpoint}?{WebService.IndexField}={HttpIndexPlaceholder}&{WebService.DocumentIdField}={HttpDocumentIdPlaceholder}"
	HttpDeleteIndexEndpointWithParams    = "{HttpIndexesEndpoint}?{WebService.IndexField}={HttpIndexPlaceholder}"
	HttpUploadStatusEndpointWithParams   = "{HttpUploadStatusEndpoint}?{WebService.IndexField}={HttpIndexPlaceholder}&{WebService.DocumentIdField}={HttpDocumentIdPlaceholder}"
	HttpDownloadEndpointWithParams       = "{HttpDownloadEndpoint}?{WebService.IndexField}={HttpIndexPlaceholder}&{WebService.DocumentIdField}={HttpDocumentIdPlaceholder}&{WebService.FilenameField}={HttpFilenamePlaceholder}"
)

const (
	// Partitioning
	Partitioning_MaxTokensPerChunk = "custom_partitioning_max_tokens_per_paragraph_int"
	Partitioning_OverlappingTokens = "custom_partitioning_overlapping_tokens_int"
	Partitioning_ChunkHeader       = "custom_partitioning_chunk_header_str"

	// EmbeddingGeneration
	EmbeddingGeneration_BatchSize = "custom_embedding_generation_batch_size_int"
	EmbeddingGeneration_ModelName = "custom_embedding_generation_model_name"

	// TextGeneration
	TextGeneration_ModelName = "custom_text_generation_model_name"

	// Rag
	Rag_EmptyAnswer           = "custom_rag_empty_answer_str"
	Rag_Prompt                = "custom_rag_prompt_str"
	Rag_FactTemplate          = "custom_rag_fact_template_str"
	Rag_IncludeDuplicateFacts = "custom_rag_include_duplicate_facts_bool"
	Rag_MaxTokens             = "custom_rag_max_tokens_int"
	Rag_MaxMatchesCount       = "custom_rag_max_matches_count_int"
	Rag_Temperature           = "custom_rag_temperature_float"
	Rag_NucleusSampling       = "custom_rag_nucleus_sampling_float"

	// Summary
	Summary_Prompt            = "custom_summary_prompt_str"
	Summary_TargetTokenSize   = "custom_summary_target_token_size_int"
	Summary_OverlappingTokens = "custom_summary_overlapping_tokens_int"
)

var (
	// Pipeline Steps
	PipelineStepsExtract              = "extract"
	PipelineStepsPartition            = "partition"
	PipelineStepsGenEmbeddings        = "gen_embeddings"
	PipelineStepsSaveRecords          = "save_records"
	PipelineStepsSummarize            = "summarize"
	PipelineStepsDeleteGeneratedFiles = "delete_generated_files"
	PipelineStepsDeleteDocument       = "private_delete_document"
	PipelineStepsDeleteIndex          = "private_delete_index"

	// Default Pipelines
	DefaultPipeline = []string{
		PipelineStepsExtract, PipelineStepsPartition, PipelineStepsGenEmbeddings, PipelineStepsSaveRecords,
	}

	PipelineWithoutSummary = []string{
		PipelineStepsExtract, PipelineStepsPartition, PipelineStepsGenEmbeddings, PipelineStepsSaveRecords,
	}

	PipelineWithSummary = []string{
		PipelineStepsExtract, PipelineStepsPartition, PipelineStepsGenEmbeddings, PipelineStepsSaveRecords,
		PipelineStepsSummarize, PipelineStepsGenEmbeddings, PipelineStepsSaveRecords,
	}

	PipelineOnlySummary = []string{
		PipelineStepsExtract, PipelineStepsSummarize, PipelineStepsGenEmbeddings, PipelineStepsSaveRecords,
	}

	// Prompt Names
	PromptNamesSummarize       = "summarize"
	PromptNamesAnswerWithFacts = "answer-with-facts"
)
