package ai

import (
	"encoding/json"
	"errors"
	"math"
)

// Embedding represents the embedding vector with float32 data.
type Embedding struct {
	Data []float32 `json:"-"` // Raw data, ignored in JSON serialization.
}

// NewEmbedding creates a new Embedding from a float slice.
func NewEmbedding(vector []float32) Embedding {
	return Embedding{Data: vector}
}

// NewEmbeddingFromDoubles creates a new Embedding from a double slice.
func NewEmbeddingFromDoubles(vector []float64) Embedding {
	data := make([]float32, len(vector))
	for i, v := range vector {
		data[i] = float32(v)
	}
	return Embedding{Data: data}
}

// Length returns the length of the embedding vector.
func (e *Embedding) Length() int {
	return len(e.Data)
}

// CosineSimilarity calculates the cosine similarity between two embeddings.
func (e *Embedding) CosineSimilarity(other Embedding) (float64, error) {
	if len(e.Data) != len(other.Data) {
		return 0, errors.New("embedding vectors must have the same length")
	}

	dotProduct := float32(0)
	magnitudeA := float32(0)
	magnitudeB := float32(0)

	for i := 0; i < len(e.Data); i++ {
		dotProduct += e.Data[i] * other.Data[i]
		magnitudeA += e.Data[i] * e.Data[i]
		magnitudeB += other.Data[i] * other.Data[i]
	}

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0, nil // Avoid division by zero
	}

	return float64(dotProduct) / (math.Sqrt(float64(magnitudeA)) * math.Sqrt(float64(magnitudeB))), nil
}

// Equals checks if two embeddings are equal.
func (e *Embedding) Equals(other Embedding) bool {
	if len(e.Data) != len(other.Data) {
		return false
	}
	for i := 0; i < len(e.Data); i++ {
		if e.Data[i] != other.Data[i] {
			return false
		}
	}
	return true
}

// MarshalJSON implements custom JSON serialization.
func (e *Embedding) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Data)
}

// UnmarshalJSON implements custom JSON deserialization.
func (e *Embedding) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.Data)
}
