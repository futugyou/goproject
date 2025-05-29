package protocol

import (
	"encoding/json"
	"fmt"
)

type Content struct {
	Type        string            `json:"type"`
	Text        *string           `json:"text,omitempty"`
	Data        *string           `json:"data,omitempty"`
	MimeType    *string           `json:"mimeType,omitempty"`
	Resource    IResourceContents `json:"resource,omitempty"`
	Annotations *Annotations      `json:"annotations,omitempty"`
}

func (c *Content) MarshalJSON() ([]byte, error) {
	type Alias Content
	aux := &struct {
		*Alias
		Resource interface{} `json:"resource,omitempty"`
	}{
		Alias: (*Alias)(c),
	}

	if c.Resource != nil {
		aux.Resource = c.Resource
	}

	return json.Marshal(aux)
}

func (c *Content) UnmarshalJSON(data []byte) error {
	type Alias Content
	aux := &struct {
		*Alias
		Resource json.RawMessage `json:"resource,omitempty"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.Resource != nil {
		var typeInfo struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(aux.Resource, &typeInfo); err != nil {
			return fmt.Errorf("failed to get resource type: %w", err)
		}

		factory, ok := resourceFactories[typeInfo.Type]
		if !ok {
			return fmt.Errorf("unknown type: %s", typeInfo.Type)
		}
		resource := factory()
		if err := json.Unmarshal(aux.Resource, resource); err != nil {
			return err
		}

		if resource != nil {
			c.Resource = resource
		}
	}

	return nil
}
