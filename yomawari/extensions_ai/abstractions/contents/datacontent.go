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

func NewDataContentFull(uri string, mediaType string, rawRepresentation interface{}, additionalProperties map[string]interface{}) *DataContent {
	content := &DataContent{
		AIContent: NewAIContent(rawRepresentation, additionalProperties),
		MediaType: mediaType,
	}

	var err error
	content.dataURI, err = ParseDataUri(uri)
	if err != nil {
		return content
	}

	if mediaType == "" {
		mediaType = content.dataURI.MediaType
	}

	if mediaType != "" {
		content.MediaType, err = ThrowIfInvalidMediaType(mediaType, mediaType)
		if err != nil {
			return content
		}
		if !content.dataURI.IsBase64 || content.MediaType != content.dataURI.MediaType {
			content.Data = content.dataURI.Data
			content.dataURI = nil
			content.URI = ""
		}

	}

	return content
}

func NewDataContent(uri string, mediaType string) *DataContent {
	return NewDataContentFull(uri, mediaType, nil, nil)
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

func (dc *DataContent) GetBase64Data() []byte {
	index := strings.Index(dc.URI, ",")
	if index == -1 || index+1 >= len(dc.URI) {
		return []byte{}
	}
	return []byte(dc.URI[index+1:])
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

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	newDC := NewDataContent(ac.URI, ac.MediaType)
	*ac = *newDC

	return nil
}
