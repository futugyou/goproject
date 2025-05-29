package protocol

import (
	"encoding/json"
	"errors"
)

type ProgressNotification struct {
	ProgressToken    *ProgressToken
	Progress         *ProgressNotificationValue
	NotificationType string
}

type rawProgressNotification struct {
	ProgressToken json.RawMessage `json:"progressToken"`
	Progress      float32         `json:"progress"`
	Total         *float32        `json:"total,omitempty"`
	Message       *string         `json:"message,omitempty"`
}

func (p *ProgressNotification) UnmarshalJSON(data []byte) error {
	var raw rawProgressNotification
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.ProgressToken == nil {
		return errors.New("missing required property 'progressToken'")
	}
	if raw.Progress == 0 {
		return errors.New("missing required property 'progress'")
	}

	var token ProgressToken
	if err := json.Unmarshal(raw.ProgressToken, &token); err != nil {
		return err
	}

	p.ProgressToken = &token
	p.Progress = &ProgressNotificationValue{
		Progress: raw.Progress,
		Total:    raw.Total,
		Message:  raw.Message,
	}
	return nil
}

func (p ProgressNotification) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ProgressToken ProgressToken `json:"progressToken"`
		Progress      float32       `json:"progress"`
		Total         *float32      `json:"total,omitempty"`
		Message       *string       `json:"message,omitempty"`
	}{
		ProgressToken: *p.ProgressToken,
		Progress:      p.Progress.Progress,
		Total:         p.Progress.Total,
		Message:       p.Progress.Message,
	})
}
