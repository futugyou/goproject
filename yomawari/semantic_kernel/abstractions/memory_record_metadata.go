package abstractions

type MemoryRecordMetadata struct {
	IsReference        bool   `json:"is_reference"`
	ExternalSourceName string `json:"external_source_name"`
	Id                 string `json:"id"`
	Description        string `json:"description"`
	Text               string `json:"text"`
	AdditionalMetadata string `json:"additional_metadata"`
}
