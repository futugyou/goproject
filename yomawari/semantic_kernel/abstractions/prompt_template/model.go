package prompt_template

import "github.com/futugyou/yomawari/semantic_kernel/abstractions/ai_functional"

type TemplateOptions struct {
	Format string
	Parser string
	Strict bool
}

type InputVariable struct {
	Name                       string `json:"name"`
	Description                string `json:"description"`
	Default                    any    `json:"default"`
	IsRequired                 bool   `json:"is_required"`
	JsonSchema                 string `json:"json_schema"`
	AllowDangerouslySetContent bool   `json:"allow_dangerously_set_content"`
	Sample                     any    `json:"sample"`
}

type OutputVariable struct {
	JsonSchema  string `json:"json_schema"`
	Description string `json:"description"`
}

type PromptTemplateConfig struct {
	Name                       string                                           `json:"name"`
	Description                string                                           `json:"description"`
	TemplateFormat             string                                           `json:"template_format"`
	Template                   string                                           `json:"template"`
	InputVariables             []InputVariable                                  `json:"input_variables"`
	OutputVariable             OutputVariable                                   `json:"output_variable"`
	ExecutionSettings          map[string]ai_functional.PromptExecutionSettings `json:"execution_settings"`
	AllowDangerouslySetContent bool                                             `json:"allow_dangerously_set_content"`
}
