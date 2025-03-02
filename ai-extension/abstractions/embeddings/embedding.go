package embeddings

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type IEmbedding interface {
	Hash() string
}

type Embedding struct {
	CreatedAt            *time.Time             `json:"createdAt,omitempty"`
	ModelId              *string                `json:"modelId,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

func (d Embedding) Hash() string {
	data, _ := json.Marshal(d)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (e Embedding) IsEquals(b Embedding) bool {
	if (e.CreatedAt == nil) != (b.CreatedAt == nil) {
		return false
	}
	if e.CreatedAt != nil && *e.CreatedAt != *b.CreatedAt {
		return false
	}

	if (e.ModelId == nil) != (b.ModelId == nil) {
		return false
	}
	if e.ModelId != nil && *e.ModelId != *b.ModelId {
		return false
	}

	if (e.AdditionalProperties == nil) != (b.AdditionalProperties == nil) {
		return false
	}
	if e.AdditionalProperties != nil {
		if len(e.AdditionalProperties) != len(b.AdditionalProperties) {
			return false
		}
		for key, value := range e.AdditionalProperties {
			bValue, exists := b.AdditionalProperties[key]
			if !exists || !reflect.DeepEqual(value, bValue) {
				return false
			}
		}
	}

	return true
}

type EmbeddingT[T any] struct {
	Embedding
	Vector []T `json:"vector"`
}

func (d EmbeddingT[T]) Hash() string {
	data, _ := json.Marshal(d)
	return fmt.Sprintf("%x", md5.Sum(data))
}

type EmbeddingType interface {
	GetType() string
}

func (e *EmbeddingT[T]) GetType() string {
	return fmt.Sprintf("%T", e)
}

func FromJSON(data []byte) (EmbeddingType, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	if t, exists := raw["type"].(string); exists {
		switch t {
		case "halves":
			var emb EmbeddingT[float32]
			if err := json.Unmarshal(data, &emb); err != nil {
				return nil, err
			}
			return &emb, nil
		case "floats":
			var emb EmbeddingT[float32]
			if err := json.Unmarshal(data, &emb); err != nil {
				return nil, err
			}
			return &emb, nil
		case "doubles":
			var emb EmbeddingT[float64]
			if err := json.Unmarshal(data, &emb); err != nil {
				return nil, err
			}
			return &emb, nil
		case "bytes":
			var emb EmbeddingT[byte]
			if err := json.Unmarshal(data, &emb); err != nil {
				return nil, err
			}
			return &emb, nil
		case "sbytes":
			var emb EmbeddingT[int8]
			if err := json.Unmarshal(data, &emb); err != nil {
				return nil, err
			}
			return &emb, nil
		}
	}

	return nil, fmt.Errorf("unknown type")
}
