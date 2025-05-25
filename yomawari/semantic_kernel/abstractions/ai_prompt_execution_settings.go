package abstractions

type PromptExecutionSettings struct {
	ServiceId              string                 `json:"service_id"`
	ModelId                string                 `json:"model_id"`
	ExtensionData          map[string]interface{} `json:"extension_data"`
	FunctionChoiceBehavior FunctionChoiceBehavior `json:"function_choice_behavior"`
}
