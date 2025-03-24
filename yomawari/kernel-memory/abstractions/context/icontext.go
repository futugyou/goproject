package context

type IContext interface {
	InitArgs(map[string]interface{}) IContext
	SetArgs(map[string]interface{}) IContext
	SetArg(key string, value interface{}) IContext
	GetArgs() map[string]interface{}
	ResetArgs() IContext
	TryGetArg(ctx IContext, key string) (interface{}, bool)
	GetCustomEmptyAnswerTextOrDefault(defaultValue string) string
	GetCustomRagFactTemplateOrDefault(defaultValue string) string
	GetCustomRagIncludeDuplicateFactsOrDefault(defaultValue bool) bool
	GetCustomRagPromptOrDefault(defaultValue string) string
	GetCustomRagMaxTokensOrDefault(defaultValue int64) int64
	GetCustomRagMaxMatchesCountOrDefault(defaultValue int64) int64
	GetCustomRagTemperatureOrDefault(defaultValue float64) float64
	GetCustomRagNucleusSamplingOrDefault(defaultValue float64) float64
	GetCustomSummaryPromptOrDefault(defaultValue string) string
	GetCustomSummaryTargetTokenSizeOrDefault(defaultValue int64) int64
	GetCustomSummaryOverlappingTokensOrDefault(defaultValue int64) int64
	GetCustomPartitioningMaxTokensPerChunkOrDefault(defaultValue int64) int64
	GetCustomPartitioningOverlappingTokensOrDefault(defaultValue int64) int64
	GetCustomPartitioningChunkHeaderOrDefault(defaultValue *string) *string
	GetCustomEmbeddingGenerationBatchSizeOrDefault(defaultValue int64) int64
	GetCustomTextGenerationModelNameOrDefault(defaultValue string) string
	GetCustomEmbeddingGenerationModelNameOrDefault(defaultValue string) string
}

func TryGetContextArg[T any](ctx IContext, key string) (*T, bool) {
	if ctx == nil {
		return new(T), false
	}
	if ctx.GetArgs() == nil {
		return new(T), false
	}
	if v, ok := ctx.GetArgs()[key]; ok {
		if condition, ok := v.(T); ok {
			return &condition, true
		}
	}
	return new(T), false
}

var _ IContext = &RequestContext{}
