package openai

import (
	"encoding/json"
)

// IJsonModel interface, similar to C#'s IJsonModel<T>
type IJsonModel[T any] interface {
	Write() ([]byte, error)                    // Serialization
	Create([]byte) (T, error)                  // Deserialization
	CreateFromReader(*json.Decoder) (T, error) // Deserialize from JSON parser
}

// Serialize - Serialize to byte array (similar to BinaryData)
func Serialize[T IJsonModel[T]](value T) ([]byte, error) {
	return value.Write()
}

// Deserialize - Deserialize JSON byte array
func Deserialize[T IJsonModel[T]](data []byte) (T, error) {
	var model T
	return model.Create(data)
}

// DeserializeFromReader - Deserialize from JSON parser
func DeserializeFromReader[T IJsonModel[T]](decoder *json.Decoder) (T, error) {
	var model T
	return model.CreateFromReader(decoder)
}
