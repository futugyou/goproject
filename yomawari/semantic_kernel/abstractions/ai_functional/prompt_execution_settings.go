package ai_functional

type PromptExecutionSettings struct {
	ServiceId     string                 `json:"service_id"`
	ModelId       string                 `json:"model_id"`
	ExtensionData map[string]interface{} `json:"extension_data"`
}
