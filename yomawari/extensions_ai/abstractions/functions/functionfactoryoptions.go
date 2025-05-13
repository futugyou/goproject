package functions

import "encoding/json"

// AIFunctionFactoryOptions holds options for creating an AIFunction.
type AIFunctionFactoryOptions struct {
	SerializerOptions    *json.Encoder
	Name                 string                 // function name
	Description          string                 // function description
	ParameterNames       []string               // parameter names
	JSONSchemaOptions    map[string]interface{} // JSON Schema
	AdditionalProperties map[string]interface{}
}
