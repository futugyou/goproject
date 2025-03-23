package memorystorage

import (
	"encoding/json"
	"sync"

	"github.com/futugyou/yomawari/kernel-memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/models"
)

const SchemaVersionZero string = ""
const SchemaVersion20231218A string = "20231218A"
const CurrentSchemaVersion string = SchemaVersion20231218A

// MemoryRecord
type MemoryRecord struct {
	mu      sync.Mutex
	Id      string                `json:"id"`
	Vector  *ai.Embedding         `json:"vector"`
	Tags    *models.TagCollection `json:"tags"`
	Payload map[string]any        `json:"payload"`
}

func NewMemoryRecord() *MemoryRecord {
	return &MemoryRecord{
		Id:      "",
		Vector:  &ai.Embedding{},
		Tags:    &models.TagCollection{},
		Payload: make(map[string]any),
	}
}

// UpgradeRequired determines whether an upgrade is required
func (m *MemoryRecord) UpgradeRequired() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	version, exists := m.Payload[constant.ReservedPayloadSchemaVersionField]
	if !exists || version == nil || version == "" {
		return true
	}

	return version != CurrentSchemaVersion
}

// Upgrade to the latest schema
func (m *MemoryRecord) Upgrade() {
	m.mu.Lock()
	defer m.mu.Unlock()

	version, exists := m.Payload[constant.ReservedPayloadSchemaVersionField]
	if !exists || version == nil {
		version = SchemaVersionZero
	}

	if version == SchemaVersionZero {
		if _, exists := m.Payload[constant.ReservedPayloadUrlField]; !exists {
			m.Payload[constant.ReservedPayloadUrlField] = ""
		}
		m.Payload[constant.ReservedPayloadSchemaVersionField] = SchemaVersion20231218A
	}

	m.Payload[constant.ReservedPayloadSchemaVersionField] = CurrentSchemaVersion
}

// MarshalJSON
func (m *MemoryRecord) MarshalJSON() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.UpgradeRequired() {
		m.Upgrade()
	}

	type Alias MemoryRecord
	return json.Marshal((*Alias)(m))
}

// UnmarshalJSON
func (m *MemoryRecord) UnmarshalJSON(data []byte) error {
	type Alias struct {
		Id      string                `json:"id"`
		Vector  *ai.Embedding         `json:"vector"`
		Tags    *models.TagCollection `json:"tags"`
		Payload map[string]any        `json:"payload"`
	}

	aux := &Alias{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.Id = aux.Id
	m.Vector = aux.Vector
	m.Tags = aux.Tags
	m.Payload = aux.Payload

	if m.UpgradeRequired() {
		m.Upgrade()
	}

	return nil
}
