package abstractions

import (
	"encoding/json"
	"time"
)

type MemoryRecord struct {
	Key       string               `json:"key"`
	Timestamp *time.Time           `json:"timestamp,omitempty"`
	Embedding []float32            `json:"embedding"`
	Metadata  MemoryRecordMetadata `json:"metadata"`
}

func (m MemoryRecord) HasTimestamp() bool {
	return m.Timestamp != nil
}

func (m MemoryRecord) GetSerializedMetadata() string {
	d, _ := json.Marshal(m.Metadata)
	if len(d) > 0 {
		return string(d)
	}
	return ""
}

func ReferenceRecord(externalId string, sourceName string, description *string, embedding []float32, additionalMetadata *string, key *string, timestamp *time.Time) MemoryRecord {
	m := MemoryRecord{
		Key:       "",
		Timestamp: timestamp,
		Embedding: embedding,
		Metadata: MemoryRecordMetadata{
			IsReference:        true,
			ExternalSourceName: sourceName,
			Id:                 externalId,
			Description:        "",
			Text:               "",
			AdditionalMetadata: "",
		},
	}

	if key != nil {
		m.Key = *key
	}

	if additionalMetadata != nil {
		m.Metadata.AdditionalMetadata = *additionalMetadata
	}

	if description != nil {
		m.Metadata.Description = *description
	}
	return m
}

func LocalRecord(id string, text string, description *string, embedding []float32, additionalMetadata *string, key *string, timestamp *time.Time) MemoryRecord {
	m := MemoryRecord{
		Key:       "",
		Timestamp: timestamp,
		Embedding: embedding,
		Metadata: MemoryRecordMetadata{
			IsReference:        false,
			ExternalSourceName: "",
			Id:                 id,
			Description:        "",
			Text:               text,
			AdditionalMetadata: "",
		},
	}

	if key != nil {
		m.Key = *key
	}

	if additionalMetadata != nil {
		m.Metadata.AdditionalMetadata = *additionalMetadata
	}

	if description != nil {
		m.Metadata.Description = *description
	}
	return m
}

func FromJsonMetadata(json_string string, embedding []float32, key *string, timestamp *time.Time) MemoryRecord {
	var m MemoryRecordMetadata
	json.Unmarshal([]byte(json_string), &m)
	e := MemoryRecord{
		Key:       "",
		Timestamp: timestamp,
		Embedding: embedding,
		Metadata:  m,
	}
	if key != nil {
		e.Key = *key
	}
	return e
}

func FromMetadata(metadata MemoryRecordMetadata, embedding []float32, key *string, timestamp *time.Time) MemoryRecord {
	e := MemoryRecord{
		Key:       "",
		Timestamp: timestamp,
		Embedding: embedding,
		Metadata:  metadata,
	}
	if key != nil {
		e.Key = *key
	}
	return e
}
