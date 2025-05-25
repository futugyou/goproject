package abstractions

import (
	"encoding/json"
	"fmt"
)

type KernelContent interface {
	Type() string
	ToString() string
	GetInnerContent() any
	Hash() string
}

func MarshalKernelContent(con KernelContent) ([]byte, error) {
	var payload map[string]any
	contentBytes, err := json.Marshal(con)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(contentBytes, &payload)
	payload["type"] = con.Type()
	return json.Marshal(payload)
}

func UnmarshalKernelContent(data []byte) (KernelContent, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	t, ok := raw["type"]
	if !ok {
		return nil, fmt.Errorf("missing type field")
	}

	var typeStr string
	if err := json.Unmarshal(t, &typeStr); err != nil {
		return nil, err
	}

	return DecodeContent(data, typeStr)
}

var (
	contentRegistry = map[string]func() KernelContent{}
)

func RegisterContent(name string, ctor func() KernelContent) {
	contentRegistry[name] = ctor
}

func DecodeContent(raw json.RawMessage, t string) (KernelContent, error) {
	ctor, ok := contentRegistry[t]
	if !ok {
		return nil, fmt.Errorf("unknown content type: %s", t)
	}
	instance := ctor()
	if err := json.Unmarshal(raw, instance); err != nil {
		return nil, err
	}
	return instance, nil
}

type ContentWrapper struct {
	Content KernelContent
}

func (w ContentWrapper) MarshalJSON() ([]byte, error) {
	var payload map[string]any
	contentBytes, err := json.Marshal(w.Content)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(contentBytes, &payload)
	payload["type"] = w.Content.Type()
	return json.Marshal(payload)
}

func (w *ContentWrapper) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t, ok := raw["type"]
	if !ok {
		return fmt.Errorf("missing type field")
	}

	var typeStr string
	if err := json.Unmarshal(t, &typeStr); err != nil {
		return err
	}

	content, err := DecodeContent(data, typeStr)
	if err != nil {
		return err
	}

	w.Content = content
	return nil
}
