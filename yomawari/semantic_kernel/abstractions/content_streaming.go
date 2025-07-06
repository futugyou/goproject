package abstractions

import (
	"encoding/base64"
	"encoding/json"
)

var _ StreamingKernelContent = (*StreamingMethodContent)(nil)

type StreamingMethodContent struct {
	Content  any
	Metadata map[string]any `json:"metadata"`
}

// Hash implements StreamingKernelContent.
func (s *StreamingMethodContent) Hash() string {
	return s.ToString()
}

// ToByteArray implements StreamingKernelContent.
func (s *StreamingMethodContent) ToByteArray() []byte {
	if d, ok := s.Content.([]byte); ok {
		return d
	}
	r, err := base64.URLEncoding.DecodeString(s.ToString())
	if err != nil {
		return []byte{}
	}
	return r
}

// ToString implements StreamingKernelContent.
func (s *StreamingMethodContent) ToString() string {
	d, err := json.Marshal(s.Content)
	if err != nil {
		return ""
	}
	return string(d)
}

// Type implements StreamingKernelContent.
func (s *StreamingMethodContent) Type() string {
	return "StreamingMethodContent"
}
