package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// Test structure to implement IJsonModel interface
type TestModel struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

var ErrTestError = errors.New("test error")

func (t TestModel) Write() ([]byte, error) {
	if t.Name == "error" {
		return nil, ErrTestError
	}
	return json.Marshal(t)
}

func (t TestModel) Create(data []byte) (TestModel, error) {
	if len(data) == 0 {
		return TestModel{}, ErrTestError
	}

	var m TestModel
	err := json.Unmarshal(data, &m)
	if err != nil {
		return TestModel{}, err
	}
	return m, nil
}

func (t TestModel) CreateFromReader(decoder *json.Decoder) (TestModel, error) {
	if decoder == nil {
		return TestModel{}, ErrTestError
	}

	var m TestModel
	err := decoder.Decode(&m)
	if err != nil {
		return TestModel{}, err
	}
	return m, nil
}

func TestSerialize(t *testing.T) {
	t.Run("Normal serialization", func(t *testing.T) {
		model := TestModel{Name: "test", Value: 42}
		data, err := Serialize(model)
		if err != nil {
			t.Fatalf("No errors are expected, and we got: %v", err)
		}

		expected := `{"name":"test","value":42}`
		if string(data) != expected {
			t.Errorf("expected %s, got %s", expected, string(data))
		}
	})

	t.Run("serialization error", func(t *testing.T) {
		model := TestModel{Name: "error"}
		_, err := Serialize(model)
		if !errors.Is(err, ErrTestError) {
			t.Errorf("Expected Errors %v, got %v", ErrTestError, err)
		}
	})
}

func TestDeserialize(t *testing.T) {
	t.Run("Normal deserialization", func(t *testing.T) {
		data := []byte(`{"name":"test","value":42}`)
		result, err := Deserialize[TestModel](data)
		if err != nil {
			t.Fatalf("Expected no errors, got: %v", err)
		}

		if result.Name != "test" || result.Value != 42 {
			t.Errorf("The deserialization result does not match, got: %+v", result)
		}
	})

	t.Run("Empty data error handling", func(t *testing.T) {
		_, err := Deserialize[TestModel]([]byte{})
		if !errors.Is(err, ErrTestError) {
			t.Errorf("Expected Errors %v, got %v", ErrTestError, err)
		}
	})

	t.Run("Invalid JSON processing", func(t *testing.T) {
		data := []byte(`{invalid json}`)
		_, err := Deserialize[TestModel](data)
		if err == nil {
			t.Error("Expected JSON parsing error, but no error was returned")
		}
	})
}

func TestDeserializeFromReader(t *testing.T) {
	t.Run("Normal streaming deserialization", func(t *testing.T) {
		jsonStr := `{"name":"stream","value":99}`
		decoder := json.NewDecoder(strings.NewReader(jsonStr))
		result, err := DeserializeFromReader[TestModel](decoder)
		if err != nil {
			t.Fatalf("Expected no errors, got: %v", err)
		}

		if result.Name != "stream" || result.Value != 99 {
			t.Errorf("The deserialization result does not match, got: %+v", result)
		}
	})

	t.Run("Empty decoder handling", func(t *testing.T) {
		_, err := DeserializeFromReader[TestModel](nil)
		if !errors.Is(err, ErrTestError) {
			t.Errorf("Expected Errors %v, got %v", ErrTestError, err)
		}
	})

	t.Run("Decoding error handling", func(t *testing.T) {
		jsonStr := `{"name": "invalid", "value": "string"}`
		decoder := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
		_, err := DeserializeFromReader[TestModel](decoder)
		if err == nil {
			t.Error("A type conversion error was expected, but no error was returned")
		}
	})
}

func TestInterfaceImplementation(t *testing.T) {
	// Verify interface implementation
	var _ IJsonModel[TestModel] = TestModel{}
}
