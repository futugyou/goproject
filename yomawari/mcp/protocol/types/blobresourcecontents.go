package types

import "encoding/json"

type BlobResourceContents struct {
	BaseResourceContents `json:",inline"`
	Blob                 string `json:"blob"`
}

func (b BlobResourceContents) MarshalJSON() ([]byte, error) {
	type Alias BlobResourceContents
	return json.Marshal(&struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(b),
		Type:  "blob",
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
