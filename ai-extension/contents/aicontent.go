package contents

import (
	"fmt"
	"reflect"
)

type AIContent struct {
	RawRepresentation    interface{}            // Raw representation of the content (for debugging or underlying object model).
	AdditionalProperties map[string]interface{} // Additional properties for the content.
}

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
