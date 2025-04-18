package contents

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// DataContent represents content that holds data, similar to the C# DataContent class.
type DataContent struct {
	*AIContent `json:",inline"`
	URI        string   `json:"url,omitempty"`
	MediaType  string   `json:"mediaType,omitempty"`
	Data       []byte   `json:"-"`
	dataURI    *DataUri `json:"-"`
}

// NewDataContent creates a new DataContent from a URI string, optionally providing a media type.
func NewDataContent(uri string, mediaType string) *DataContent {
	content := &DataContent{
		AIContent: NewAIContent(nil, nil),
		MediaType: mediaType}

	// Check if URI is a data URI
	if strings.HasPrefix(uri, "data:") {
		dataURI, err := ParseDataUri(uri)
		if err != nil {
			return content
		}

		content.dataURI = dataURI
		if mediaType == "" {
			content.MediaType = dataURI.MediaType
		} else if mediaType != dataURI.MediaType {
			content.Data = dataURI.Data
			content.dataURI = nil
			content.URI = ""
		}
	} else {
		// If not a data URI, treat it as a regular URI
		content.URI = uri
	}

	return content
}

func NewDataContentWithRefusal(uri string, mediaType string, refusal string) *DataContent {
	c := NewDataContent(uri, mediaType)
	if len(refusal) > 0 {
		c.AdditionalProperties["refusal"] = refusal
	}
	return c
}

// GetURI returns the URI for this DataContent.
func (dc *DataContent) GetURI() string {
	if dc.URI == "" && dc.dataURI != nil {
		if dc.dataURI.IsBase64 {
			dc.URI = fmt.Sprintf("data:%s;base64,%s", dc.MediaType, base64.StdEncoding.EncodeToString(dc.dataURI.Data))
		} else {
			dc.URI = fmt.Sprintf("data:%s,%s", dc.MediaType, string(dc.dataURI.Data))
		}
	}
	return dc.URI
}

// GetData returns the data for this DataContent, if available.
func (dc *DataContent) GetData() ([]byte, error) {
	if dc.dataURI != nil && len(dc.Data) == 0 {
		dc.Data = dc.dataURI.Data
	}
	if len(dc.Data) > 0 {
		return dc.Data, nil
	}
	return nil, errors.New("no data available")
}

func (dc *DataContent) MediaTypeStartsWith(prefix string) bool {
	return strings.HasPrefix(dc.MediaType, prefix)
}

func (ac DataContent) MarshalJSON() ([]byte, error) {
	type Alias DataContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "DataContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *DataContent) UnmarshalJSON(data []byte) error {
	type Alias DataContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
