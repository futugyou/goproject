package memorystorage

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
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
	type Alias MemoryRecord
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

func (m *MemoryRecord) GetTagValue(tagName string) string {
	if tagValues, ok := m.Tags.Get(tagName); ok {
		if len(tagValues) > 0 {
			return tagValues[0]
		}
	}

	return ""
}

func (m *MemoryRecord) GetDocumentId() string {
	return m.GetTagValue(constant.ReservedDocumentIdTag)
}

func (m *MemoryRecord) GetFileId() string {
	return m.GetTagValue(constant.ReservedFileIdTag)
}

func (m *MemoryRecord) GetPartitionNumber() int {
	value := m.GetTagValue(constant.ReservedFilePartitionNumberTag)
	if v, err := strconv.Atoi(value); err != nil {
		return v
	}
	return 0
}

func (m *MemoryRecord) GetSectionNumber() int {
	value := m.GetTagValue(constant.ReservedFileSectionNumberTag)
	if v, err := strconv.Atoi(value); err != nil {
		return v
	}
	return 0
}

func (m *MemoryRecord) GetFileContentType() string {
	return m.GetTagValue(constant.ReservedFileTypeTag)
}

func (m *MemoryRecord) GetFileName() string {
	if v, ok := m.GetPayloadValue(constant.ReservedPayloadFileNameField).(string); ok {
		return v
	}
	return ""
}

func (m *MemoryRecord) GetWebPageUrl(indexName string) string {
	if webPageUrl, ok := m.GetPayloadValue(constant.ReservedPayloadUrlField).(string); ok && len(webPageUrl) > 0 {
		return webPageUrl
	}
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				constant.HttpDownloadEndpointWithParams,
				constant.HttpIndexPlaceholder,
				indexName,
			),
			constant.HttpDocumentIdPlaceholder,
			m.GetDocumentId(),
		),
		constant.HttpFilenamePlaceholder,
		m.GetFileName(),
	)
}

func (m *MemoryRecord) GetPartitionText() string {
	if v, ok := m.GetPayloadValue(constant.ReservedPayloadTextField).(string); ok {
		return v
	}
	return ""
}

func (m *MemoryRecord) GetLastUpdate() time.Time {
	if v, ok := m.GetPayloadValue(constant.ReservedPayloadLastUpdateField).(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
	}

	return time.Time{}
}

func (m *MemoryRecord) GetPayloadValue(payloadKey string) any {
	if tagValues, ok := m.Payload[payloadKey]; ok {
		return tagValues
	}

	return nil
}
