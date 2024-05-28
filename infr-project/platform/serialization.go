package platform

import (
	"encoding/json"
	"errors"
)

// MarshalJSON is a custom marshaler for Webhook that handles the serialization of WebhookState.
// In this case, we can skip MarshalJSON, only implement UnmarshalJSON
func (w Webhook) MarshalJSON() ([]byte, error) {
	type Alias Webhook
	return json.Marshal(&struct {
		State string `json:"state"`
		*Alias
	}{
		State: w.State.String(),
		Alias: (*Alias)(&w),
	})
}

// UnmarshalJSON is a custom unmarshaler for Webhook that handles the deserialization of WebhookState.
func (w *Webhook) UnmarshalJSON(data []byte) error {
	type Alias Webhook
	aux := &struct {
		State string `json:"state"`
		*Alias
	}{
		Alias: (*Alias)(w),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch aux.State {
	case string(WebhookInit):
		w.State = WebhookInit
	case string(WebhookCreating):
		w.State = WebhookCreating
	case string(WebhookReady):
		w.State = WebhookReady
	default:
		return errors.New("invalid webhook state")
	}
	return nil
}
