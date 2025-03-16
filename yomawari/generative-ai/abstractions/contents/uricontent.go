package contents

type UriContent struct {
	AIContent `json:",inline"`
	URI       string `json:"url,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}
