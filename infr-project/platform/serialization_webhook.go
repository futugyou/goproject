package platform

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// MarshalJSON is a custom marshaler for Webhook that handles the serialization of WebhookState.
// In this case, we can skip MarshalJSON, only implement UnmarshalJSON
func (w Webhook) MarshalJSON() ([]byte, error) {
	return w.commonMarshal(json.Marshal)
}

func (w Webhook) MarshalBSON() ([]byte, error) {
	return w.commonMarshal(bson.Marshal)
}

func (r Webhook) commonMarshal(marshal func(interface{}) ([]byte, error)) ([]byte, error) {
	m := map[string]interface{}{
		"name":       r.Name,
		"url":        r.Url,
		"activate":   r.Activate,
		"properties": r.Properties,
	}
	if r.State != nil {
		m["state"] = r.State.String()
	}
	return marshal(m)
}

// UnmarshalJSON is a custom unmarshaler for Webhook that handles the deserialization of WebhookState.
func (w *Webhook) UnmarshalJSON(data []byte) error {
	return w.commonUnmarshal(data, json.Unmarshal)
}

func (w *Webhook) UnmarshalBSON(data []byte) error {
	return w.commonUnmarshal(data, bson.Unmarshal)
}

func (w *Webhook) commonUnmarshal(data []byte, unmarshal func([]byte, any) error) error {
	var m map[string]interface{}
	if err := unmarshal(data, &m); err != nil {
		return err
	}

	if value, ok := m["name"].(string); ok {
		w.Name = value
	}

	if value, ok := m["url"].(string); ok {
		w.Url = value
	}

	if value, ok := m["activate"].(bool); ok {
		w.Activate = value
	}

	if value, ok := m["properties"].(map[string]interface{}); ok {
		properties := make(map[string]string, len(value))
		for key, v := range value {
			if d, ok := v.(string); ok {
				properties[key] = d
			}
		}
		w.Properties = properties
	}

	if state, ok := m["state"].(string); ok {
		switch state {
		case string(WebhookInit):
			w.State = WebhookInit
		case string(WebhookCreating):
			w.State = WebhookCreating
		case string(WebhookReady):
			w.State = WebhookReady
		default:
			return errors.New("invalid webhook state")
		}
	}

	return nil
}
