package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ProgressToken struct {
	id any
}

func NewProgressTokenFromString(value string) ProgressToken {
	return ProgressToken{id: value}
}

func NewProgressTokenFromInt(value int64) ProgressToken {
	return ProgressToken{id: value}
}

func (r ProgressToken) IsDefault() bool {
	return r.id == nil
}

func (r ProgressToken) String() string {
	switch v := r.id.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

func (r ProgressToken) MarshalJSON() ([]byte, error) {
	switch v := r.id.(type) {
	case string:
		return json.Marshal(v)
	case int64:
		return json.Marshal(v)
	case nil:
		return json.Marshal("")
	default:
		return nil, errors.New("invalid ProgressToken type")
	}
}

func (r *ProgressToken) UnmarshalJSON(data []byte) error {
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

	return errors.New("ProgressToken must be a string or an integer")
}
