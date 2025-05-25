package abstractions

import (
	"bytes"
	"encoding/json"
)

// schemaStr := `{"type":"object","properties":{"name":{"type":"string"}}}`
//
// schema, err := Parse(schemaStr)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// fmt.Println(schema.String())
type KernelJsonSchema struct {
	rootElement    json.RawMessage
	schemaAsString *string
}

// ParseOrNull parses a JSON schema or returns nil if empty.
func KernelJsonSchemaParseOrNull(jsonSchema string) (*KernelJsonSchema, error) {
	if jsonSchema == "" {
		return nil, nil
	}
	return KernelJsonSchemaParse(jsonSchema)
}

// Parse parses a JSON schema string and returns a KernelJsonSchema.
func KernelJsonSchemaParse(jsonSchema string) (*KernelJsonSchema, error) {
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(jsonSchema), &raw); err != nil {
		return nil, err
	}
	return &KernelJsonSchema{rootElement: raw}, nil
}

// ParseFromBytes parses from UTF-8 JSON bytes.
func KernelJsonSchemaParseFromBytes(jsonSchema []byte) (*KernelJsonSchema, error) {
	var raw json.RawMessage
	if err := json.Unmarshal(jsonSchema, &raw); err != nil {
		return nil, err
	}
	return &KernelJsonSchema{rootElement: raw}, nil
}

// RootElement returns the root JSON element.
func (k *KernelJsonSchema) RootElement() json.RawMessage {
	return k.rootElement
}

// String returns the original or serialized string.
func (k *KernelJsonSchema) String() string {
	if k.schemaAsString != nil {
		return *k.schemaAsString
	}
	s := string(k.rootElement)
	k.schemaAsString = &s
	return s
}

// MarshalJSON implements custom JSON serialization.
func (k KernelJsonSchema) MarshalJSON() ([]byte, error) {
	return k.rootElement, nil
}

// UnmarshalJSON implements custom JSON deserialization.
func (k *KernelJsonSchema) UnmarshalJSON(data []byte) error {
	// Validate JSON
	var raw json.RawMessage
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // optional strict mode
	if err := dec.Decode(&raw); err != nil {
		return err
	}
	k.rootElement = raw
	k.schemaAsString = nil
	return nil
}
