package utilities

// AIJsonSchemaCreateOptions defines configuration options when generating JSON Schema
type AIJsonSchemaCreateOptions struct {
	IncludeSchemaKeyword         bool
	IncludeTypeInEnumSchemas     bool
	DisallowAdditionalProperties bool
	RequireAllProperties         bool
	TransformOptions             *AIJsonSchemaTransformOptions
}

// DefaultAIJsonSchemaCreateOptions default options
var DefaultAIJsonSchemaCreateOptions = &AIJsonSchemaCreateOptions{
	IncludeSchemaKeyword:         true,
	IncludeTypeInEnumSchemas:     false,
	DisallowAdditionalProperties: true,
	RequireAllProperties:         true,
}
