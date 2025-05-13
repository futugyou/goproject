package context

import "github.com/futugyou/yomawari/kernel_memory/abstractions/constant"

type RequestContext struct {
	Arguments map[string]interface{} `json:"args"`
}

func NewRequestContext(args map[string]interface{}) *RequestContext {
	return &RequestContext{Arguments: args}
}

func (rc *RequestContext) InitArgs(args map[string]interface{}) IContext {
	if rc == nil {
		return nil
	}
	rc.Arguments = args
	return rc
}

func (rc *RequestContext) SetArgs(args map[string]interface{}) IContext {
	if rc == nil {
		return nil
	}

	for key, v := range args {
		rc.Arguments[key] = v
	}
	return rc
}

func (rc *RequestContext) SetArg(key string, value interface{}) IContext {
	if rc == nil {
		return nil
	}

	if rc.Arguments == nil {
		rc.Arguments = make(map[string]interface{})
	}

	rc.Arguments[key] = value

	return rc
}

func (rc *RequestContext) GetArgs() map[string]interface{} {
	if rc == nil {
		return nil
	}
	return rc.Arguments
}

func (rc *RequestContext) ResetArgs() IContext {
	if rc == nil {
		return nil
	}
	rc.Arguments = make(map[string]interface{})
	return rc
}

func (rc *RequestContext) TryGetArg(ctx IContext, key string) (interface{}, bool) {
	if rc == nil {
		return nil, false
	}

	v, ok := rc.Arguments[key]
	return v, ok
}

func (rc *RequestContext) GetCustomEmptyAnswerTextOrDefault(defaultValue string) string {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_EmptyAnswer].(string); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomRagFactTemplateOrDefault(defaultValue string) string {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_FactTemplate].(string); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomRagIncludeDuplicateFactsOrDefault(defaultValue bool) bool {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_IncludeDuplicateFacts].(bool); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomRagPromptOrDefault(defaultValue string) string {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_Prompt].(string); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomRagMaxTokensOrDefault(defaultValue int64) int64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_MaxTokens].(int64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomRagMaxMatchesCountOrDefault(defaultValue int64) int64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_MaxMatchesCount].(int64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomRagTemperatureOrDefault(defaultValue float64) float64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_Temperature].(float64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomRagNucleusSamplingOrDefault(defaultValue float64) float64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Rag_NucleusSampling].(float64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomSummaryPromptOrDefault(defaultValue string) string {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Summary_Prompt].(string); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomSummaryTargetTokenSizeOrDefault(defaultValue int64) int64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Summary_TargetTokenSize].(int64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomSummaryOverlappingTokensOrDefault(defaultValue int64) int64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Summary_OverlappingTokens].(int64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomPartitioningMaxTokensPerChunkOrDefault(defaultValue int64) int64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Partitioning_MaxTokensPerChunk].(int64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomPartitioningOverlappingTokensOrDefault(defaultValue int64) int64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Partitioning_OverlappingTokens].(int64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomPartitioningChunkHeaderOrDefault(defaultValue *string) *string {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.Partitioning_ChunkHeader].(*string); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomEmbeddingGenerationBatchSizeOrDefault(defaultValue int64) int64 {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.EmbeddingGeneration_BatchSize].(int64); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomTextGenerationModelNameOrDefault(defaultValue string) string {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.TextGeneration_ModelName].(string); ok {
		return v
	}
	return defaultValue
}

func (rc *RequestContext) GetCustomEmbeddingGenerationModelNameOrDefault(defaultValue string) string {
	if rc == nil {
		return defaultValue
	}
	if v, ok := rc.Arguments[constant.EmbeddingGeneration_ModelName].(string); ok {
		return v
	}
	return defaultValue
}
