package utilities

type AIJsonSchemaTransformOptions struct {
	TransformSchemaNode             func(AIJsonSchemaTransformContext, map[string]interface{}) map[string]interface{}
	ConvertBooleanSchemas           bool
	DisallowAdditionalProperties    bool
	RequireAllProperties            bool
	UseNullableKeyword              bool
	MoveDefaultKeywordToDescription bool
}

type AIJsonSchemaTransformContext struct {
	path []string
}

func NewAIJsonSchemaTransformContext(path []string) AIJsonSchemaTransformContext {
	return AIJsonSchemaTransformContext{path: path}
}

func (ctx AIJsonSchemaTransformContext) Path() []string {
	return ctx.path
}

func (ctx AIJsonSchemaTransformContext) PropertyName() *string {
	path := ctx.path
	if len(path) >= 3 &&
		path[len(path)-2] == "properties" {
		name := path[len(path)-1]
		return &name
	}
	return nil
}

func (ctx AIJsonSchemaTransformContext) IsCollectionElementSchema() bool {
	path := ctx.path
	return len(path) >= 1 && path[len(path)-1] == "items"
}

func (ctx AIJsonSchemaTransformContext) IsDictionaryValueSchema() bool {
	path := ctx.path
	return len(path) >= 1 && path[len(path)-1] == "additionalProperties"
}
