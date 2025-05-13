package contents

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// DataUri represents the parsed Data URI.
// data:text/plain;base64,SGVsbG8gd29ybGQh
// Media Type: text/plain
// Is Base64: true
// Decoded Data: Hello world!
type DataUri struct {
	MediaType string
	Data      []byte
	IsBase64  bool
}

// ParseDataUri parses a data URI based on RFC 2397.
func ParseDataUri(dataUri string) (*DataUri, error) {
	// Ensure the data URI starts with "data:"
	if !strings.HasPrefix(dataUri, "data:") {
		return nil, fmt.Errorf("invalid data URI format: must start with 'data:'")
	}

	// Remove the "data:" scheme prefix
	dataUri = dataUri[5:]

	// Find the comma separating the metadata from the data
	commaPos := strings.Index(dataUri, ",")
	if commaPos == -1 {
		return nil, fmt.Errorf("invalid data URI format: must contain a comma separating metadata and data")
	}

	// Extract the metadata and data parts
	metadata := dataUri[:commaPos]
	data := dataUri[commaPos+1:]

	// Check if it's Base64 encoded
	isBase64 := false
	if strings.HasSuffix(metadata, ";base64") {
		metadata = metadata[:len(metadata)-7] // Remove ";base64"
		isBase64 = true
	}

	// Validate the media type
	mediaType := ""
	if !isValidMediaType(metadata, &mediaType) {
		return nil, fmt.Errorf("invalid media type in data URI")
	}

	// Decode the data
	var decodedData []byte
	var err error
	if isBase64 {
		decodedData, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, fmt.Errorf("invalid base64 data: %v", err)
		}
	} else {
		// url.QueryUnescape returns string, so convert it to []byte
		decodedString, err := url.QueryUnescape(data)
		if err != nil {
			return nil, fmt.Errorf("invalid URL-encoded data: %v", err)
		}
		decodedData = []byte(decodedString)
	}

	return &DataUri{
		MediaType: mediaType,
		Data:      decodedData,
		IsBase64:  isBase64,
	}, nil
}

// isValidMediaType validates the media type.
func isValidMediaType(metadata string, mediaType *string) bool {
	// For simplicity, validate against known common media types.
	knownTypes := map[string]bool{
		"application/json":         true,
		"application/octet-stream": true,
		"application/pdf":          true,
		"application/xml":          true,
		"image/jpeg":               true,
		"image/png":                true,
		"text/plain":               true,
		"text/html":                true,
	}

	// Trim any whitespace and validate
	metadata = strings.TrimSpace(metadata)
	if _, ok := knownTypes[metadata]; ok {
		*mediaType = metadata
		return true
	}

	// Otherwise, check if it's a valid media type (basic regex match)
	re := regexp.MustCompile(`^[a-zA-Z0-9!#$&'*+.^_` + "`" + `|~-]+/[a-zA-Z0-9!#$&'*+.^_` + "`" + `|~-]+$`)
	if re.MatchString(metadata) {
		*mediaType = metadata
		return true
	}

	return false
}

func ThrowIfInvalidMediaType(mediaType string, parameterName string) (string, error) {
	if len(mediaType) == 0 || len(parameterName) == 0 {
		return "", fmt.Errorf("invalid parameter")
	}

	if !isValidMediaType(mediaType, &mediaType) {
		return "", fmt.Errorf("an invalid media type was specified: '%s'", parameterName)
	}

	return mediaType, nil
}

// ToByteArray converts the DataUri to a byte array.
func (d *DataUri) ToByteArray() []byte {
	return d.Data
}
