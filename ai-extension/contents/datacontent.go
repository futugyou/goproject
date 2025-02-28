package contents

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// DataContent represents content that holds data, similar to the C# DataContent class.
type DataContent struct {
	AIContent
	URI       string
	MediaType string
	data      []byte
	dataURI   *DataUri
}

// NewDataContentFromURI creates a new DataContent from a URI string, optionally providing a media type.
func NewDataContentFromURI(uri string, mediaType string) (*DataContent, error) {
	if uri == "" {
		return nil, errors.New("URI cannot be empty")
	}

	content := &DataContent{MediaType: mediaType}

	// Check if URI is a data URI
	if strings.HasPrefix(uri, "data:") {
		dataURI, err := ParseDataUri(uri)
		if err != nil {
			return nil, err
		}

		content.dataURI = dataURI
		if mediaType == "" {
			content.MediaType = dataURI.MediaType
		} else if mediaType != dataURI.MediaType {
			content.data = dataURI.Data
			content.dataURI = nil
			content.URI = ""
		}

	} else {
		// If not a data URI, treat it as a regular URI
		content.URI = uri
	}

	return content, nil
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
	if dc.dataURI != nil && len(dc.data) == 0 {
		dc.data = dc.dataURI.Data
	}
	if len(dc.data) > 0 {
		return dc.data, nil
	}
	return nil, errors.New("no data available")
}
