package functions

type KernelFunctionSchemaModel struct {
	Type        string                            `json:"type"`
	Description string                            `json:"condition"`
	Properties  map[string]map[string]interface{} `json:"properties"`
	Required    []string                          `json:"required"`
}
