package protocol

import (
	"encoding/json"
	"fmt"
)

type ReadResourceResult struct {
	Contents []IResourceContents `json:"contents"`
}

func (r *ReadResourceResult) MarshalJSON() ([]byte, error) {
	contents := make([]interface{}, len(r.Contents))
	for i, item := range r.Contents {
		contents[i] = item
	}

	return json.Marshal(struct {
		Contents []interface{} `json:"contents"`
	}{
		Contents: contents,
	})
}

func (r *ReadResourceResult) UnmarshalJSON(data []byte) error {
	var raw struct {
		Contents []json.RawMessage `json:"contents"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Contents = make([]IResourceContents, len(raw.Contents))
	for i, itemData := range raw.Contents {
		var typeInfo struct{ Type string }
		if err := json.Unmarshal(itemData, &typeInfo); err != nil {
			return err
		}

		factory, ok := resourceFactories[typeInfo.Type]
		if !ok {
			return fmt.Errorf("unknown type: %s", typeInfo.Type)
		}
		resource := factory()
		if err := json.Unmarshal(itemData, resource); err != nil {
			return err
		}
		r.Contents[i] = resource
	}

	return nil
}
