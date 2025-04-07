package contents

import "strings"

type UriContent struct {
	AIContent `json:",inline"`
	URI       string `json:"url,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}

func (dc *UriContent) MediaTypeStartsWith(prefix string) bool {
	return strings.HasPrefix(dc.MediaType, prefix)
}
