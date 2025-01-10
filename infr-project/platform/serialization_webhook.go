package platform

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	properties := make([]Property, 0, len(r.Properties))
	for _, k := range r.Properties {
		properties = append(properties, k)
	}

	secrets := make([]Secret, 0, len(r.Secrets))
	for _, k := range r.Secrets {
		secrets = append(secrets, k)
	}
	m := map[string]interface{}{
		"name":       r.Name,
		"url":        r.Url,
		"activate":   r.Activate,
		"id":         r.ID,
		"properties": properties,
		"secrets":    secrets,
		"events":     r.Events,
	}
	if r.State != nil {
		m["state"] = r.State.String()
	} else {
		m["state"] = string(WebhookInit)
	}
	return marshal(m)
}

// UnmarshalJSON is a custom unmarshaler for Webhook that handles the deserialization of WebhookState.
func (w *Webhook) UnmarshalJSON(data []byte) error {
	return w.commonUnmarshal(data, json.Marshal, json.Unmarshal)
}

func (w *Webhook) UnmarshalBSON(data []byte) error {
	return w.commonUnmarshal(data, bson.Marshal, bson.Unmarshal)
}

func (w *Webhook) commonUnmarshal(data []byte, marshal func(interface{}) ([]byte, error), unmarshal func([]byte, any) error) error {
	var m map[string]interface{}
	if err := unmarshal(data, &m); err != nil {
		return err
	}

	if value, ok := m["name"].(string); ok {
		w.Name = value
	}

	if value, ok := m["id"].(string); ok {
		w.ID = value
	}

	if value, ok := m["url"].(string); ok {
		w.Url = value
	}

	if value, ok := m["activate"].(bool); ok {
		w.Activate = value
	}

	if value, ok := m["properties"].(primitive.A); ok {
		properties, err := parseArrayToMap[Property](value, marshal, unmarshal)
		if err != nil {
			return err
		}
		w.Properties = properties
	}

	if value, ok := m["secrets"].(primitive.A); ok {
		secrets, err := parseArrayToMap[Secret](value, marshal, unmarshal)
		if err != nil {
			return err
		}
		w.Secrets = secrets
	}

	state, _ := m["state"].(string)
	w.State = GetWebhookState(state)

	if value, ok := m["events"].([]string); ok {
		w.Events = value
	} else if value, ok := m["events"].(primitive.A); ok {
		events := make([]string, 0)
		for _, item := range value {
			if event, ok := item.(string); ok {
				events = append(events, event)
			}
		}
		w.Events = events
	}

	return nil
}
