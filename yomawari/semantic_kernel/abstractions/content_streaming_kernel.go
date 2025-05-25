package abstractions

import (
	"encoding/json"
	"fmt"
)

type StreamingKernelContent interface {
	Type() string
	ToString() string
	ToByteArray() []byte
	Hash() string
}

func MarshalStreamingKernelContent(con StreamingKernelContent) ([]byte, error) {
	var payload map[string]any
	contentBytes, err := json.Marshal(con)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(contentBytes, &payload)
	payload["type"] = con.Type()
	return json.Marshal(payload)
}

func UnmarshalStreamingKernelContent(data []byte) (StreamingKernelContent, error) {
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

	return DecodeStreamContent(data, typeStr)
}

var (
	streamContentRegistry = map[string]func() StreamingKernelContent{}
)

func RegisterStreamContent(name string, ctor func() StreamingKernelContent) {
	streamContentRegistry[name] = ctor
}

func DecodeStreamContent(raw json.RawMessage, t string) (StreamingKernelContent, error) {
	ctor, ok := streamContentRegistry[t]
	if !ok {
		return nil, fmt.Errorf("unknown content type: %s", t)
	}
	instance := ctor()
	if err := json.Unmarshal(raw, instance); err != nil {
		return nil, err
	}
	return instance, nil
}
