package protocol

import "encoding/json"

type BlobResourceContents struct {
	BaseResourceContents `json:",inline"`
	Blob                 string `json:"blob"`
}

func (b BlobResourceContents) MarshalJSON() ([]byte, error) {
	type Alias BlobResourceContents
	return json.Marshal(&struct {
		Alias
		ResourceType string `json:"resource_type"`
	}{
		Alias:        Alias(b),
		ResourceType: "blob",
	})
}

func (b *BlobResourceContents) UnmarshalJSON(data []byte) error {
	type Alias BlobResourceContents
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	return nil
}
