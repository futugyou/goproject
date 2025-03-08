package contents

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type IAIContent interface {
	IsAIContent()
}

type AIContent struct {
	RawRepresentation    interface{}            // Raw representation of the content (for debugging or underlying object model).
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"` // Additional properties for the content.
}

func (ac AIContent) IsAIContent() {}

// AddAdditionalProperty allows adding properties to the content.
func (ac *AIContent) AddAdditionalProperty(key string, value interface{}) {
	if ac.AdditionalProperties == nil {
		ac.AdditionalProperties = make(map[string]interface{})
	}
	ac.AdditionalProperties[key] = value
}

// PrintContentInfo prints out the details of the content type.
func PrintContentInfo(content AIContent) {
	contentType := reflect.TypeOf(content)
	fmt.Printf("Content Type: %s\n", contentType.Name())
	if len(content.AdditionalProperties) > 0 {
		fmt.Println("Additional Properties:")
		for key, value := range content.AdditionalProperties {
			fmt.Printf(" - %s: %v\n", key, value)
		}
	}
}

func (ac AIContent) MarshalJSON() ([]byte, error) {
	type Alias AIContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  "AIContent",
		Alias: Alias(ac),
	})
}

func (ac *AIContent) UnmarshalJSON(data []byte) error {
	type Alias AIContent
	aux := &struct {
		Type string `json:"type"`
		Alias
	}{Alias: Alias(*ac)}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	*ac = AIContent(aux.Alias)
	return nil
}


func ConcatTextContents(contents []IAIContent) string {
	var text string
	for _, content := range contents {
		if textContent, ok := content.(TextContent); ok {
			text += textContent.Text
		}
	}
	return text
}