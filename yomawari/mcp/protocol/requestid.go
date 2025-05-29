package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
)

type RequestId struct {
	id any
}

func NewRequestIdFromString(value string) *RequestId {
	return &RequestId{id: value}
}

func NewRequestIdFromInt(value int64) *RequestId {
	return &RequestId{id: value}
}

func (r RequestId) IsDefault() bool {
	return r.id == nil
}

func (r *RequestId) String() string {
	if r == nil {
		return ""
	}
	switch v := r.id.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

func (r RequestId) MarshalJSON() ([]byte, error) {
	switch v := r.id.(type) {
	case string:
		return json.Marshal(v)
	case int64:
		return json.Marshal(v)
	case nil:
		return json.Marshal("")
	default:
		return nil, errors.New("invalid RequestId type")
	}
}

func (r *RequestId) UnmarshalJSON(data []byte) error {
	var strValue string
	if err := json.Unmarshal(data, &strValue); err == nil {
		r.id = strValue
		return nil
	}

	var intValue int64
	if err := json.Unmarshal(data, &intValue); err == nil {
		r.id = intValue
		return nil
	}

	return errors.New("requestId must be a string or an integer")
}
