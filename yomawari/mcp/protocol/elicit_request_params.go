package protocol

import (
	"context"
	"encoding/json"
)

type ElicitRequestParams struct {
	Message       string        `json:"message"`
	RequestSchema RequestSchema `json:"requestSchema"`
}

type PrimitiveSchemaDefinition struct {
	Type        string   `json:"type"`
	Title       *string  `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	MinLength   *int     `json:"minLength,omitempty"`
	MaxLength   *int     `json:"maxLength,omitempty"`
	Format      *string  `json:"format,omitempty"`
	Default     *bool    `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	EnumNames   []string `json:"enumNames,omitempty"`
}

type RequestSchema struct {
	Type       string                               `json:"type"`
	Properties map[string]PrimitiveSchemaDefinition `json:"properties,omitempty"`
	Required   []string                             `json:"required,omitempty"`
}

type ElicitResult struct {
	Action  string          `json:"action"`
	Content json.RawMessage `json:"content"`
}

type ElicitationCapability struct{
	ElicitationHandler func(context.Context, *ElicitRequestParams)(*ElicitResult, error)`json:"-"`
}