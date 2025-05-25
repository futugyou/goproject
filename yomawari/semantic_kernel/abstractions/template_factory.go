package abstractions

import "fmt"

type IPromptTemplateFactory interface {
	TryCreate(templateConfig PromptTemplateConfig) (IPromptTemplate, bool)
}

func CreatePromptTemplate(factory IPromptTemplateFactory, templateConfig PromptTemplateConfig) (IPromptTemplate, error) {
	result, ok := factory.TryCreate(templateConfig)
	if !ok {
		return nil, fmt.Errorf("prompt template format %s is not supported", templateConfig.TemplateFormat)
	}

	return result, nil
}
