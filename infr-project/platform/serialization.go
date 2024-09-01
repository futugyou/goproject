package platform

import (
	"encoding/json"
	"errors"
	"fmt"

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
	m := map[string]interface{}{
		"name":     r.Name,
		"url":      r.Url,
		"activate": r.Activate,
		"property": r.Property,
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

	if value, ok := m["property"].(map[string]interface{}); ok {
		property := make(map[string]string, len(value))
		for key, v := range value {
			if d, ok := v.(string); ok {
				property[key] = d
			}
		}
		w.Property = property
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

func (r Platform) MarshalJSON() ([]byte, error) {
	return json.Marshal(makeMap(r))
}

func (r *Platform) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return makeEntity(r, m, json.Marshal, json.Unmarshal)
}

func (r Platform) MarshalBSON() ([]byte, error) {
	return bson.Marshal(makeMap(r))
}

func (r *Platform) UnmarshalBSON(data []byte) error {
	var m map[string]interface{}
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	return makeEntity(r, m, bson.Marshal, bson.Unmarshal)
}

func makeEntity(r *Platform, m map[string]interface{}, marshal func(interface{}) ([]byte, error), unmarshal func([]byte, any) error) error {
	if value, ok := m["id"].(string); ok {
		r.Id = value
	}

	if value, ok := m["name"].(string); ok {
		r.Name = value
	}
	if value, ok := m["activate"].(bool); ok {
		r.Activate = value
	}

	if value, ok := m["url"].(string); ok {
		r.Url = value
	}

	if value, ok := m["tags"].([]string); ok {
		r.Tags = value
	} else if value, ok := m["tags"].(primitive.A); ok {
		tags := make([]string, 0)
		for _, item := range value {
			if tag, ok := item.(string); ok {
				tags = append(tags, tag)
			}
		}
		r.Tags = tags
	}

	if value, ok := m["is_deleted"].(bool); ok {
		r.IsDeleted = value
	}

	if value, ok := m["restendpoint"].(string); ok {
		r.RestEndpoint = value
	}

	if value, ok := m["property"].(primitive.A); ok {
		propertys := make(map[string]PropertyInfo)
		for _, item := range value {
			jsonBytes, err := marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal PropertyInfo item: %v", err)
			}

			var proinfo PropertyInfo
			if err := unmarshal(jsonBytes, &proinfo); err != nil {
				return fmt.Errorf("failed to unmarshal item to PropertyInfo: %v", err)
			}

			propertys[proinfo.Key] = proinfo
		}
		r.Property = propertys
	}

	if value, ok := m["projects"].(primitive.A); ok {
		projects := make(map[string]PlatformProject)
		for _, item := range value {
			jsonBytes, err := marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal PlatformProject item: %v", err)
			}

			var project PlatformProject
			if err := unmarshal(jsonBytes, &project); err != nil {
				return fmt.Errorf("failed to unmarshal item to PlatformProject: %v", err)
			}

			projects[project.Id] = project
		}
		r.Projects = projects
	}

	return nil
}

func makeMap(r Platform) map[string]interface{} {
	projects := make([]PlatformProject, 0, len(r.Projects))
	for _, k := range r.Projects {
		projects = append(projects, k)
	}

	properties := make([]PropertyInfo, 0, len(r.Property))
	for _, k := range r.Property {
		properties = append(properties, k)
	}

	m := map[string]interface{}{
		"id":           r.Id,
		"name":         r.Name,
		"activate":     r.Activate,
		"url":          r.Url,
		"restendpoint": r.RestEndpoint,
		"property":     properties,
		"projects":     projects,
		"tags":         r.Tags,
		"is_deleted":   r.IsDeleted,
	}

	return m
}

// PlatformProject
func (r *PlatformProject) MarshalJSON() ([]byte, error) {
	return json.Marshal(makePlatformProjectMap(r))
}

func (r *PlatformProject) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return makePlatformProjectEntity(r, m, json.Marshal, json.Unmarshal)
}

func (r *PlatformProject) MarshalBSON() ([]byte, error) {
	return bson.Marshal(makePlatformProjectMap(r))
}

func (r *PlatformProject) UnmarshalBSON(data []byte) error {
	var m map[string]interface{}
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	return makePlatformProjectEntity(r, m, bson.Marshal, bson.Unmarshal)
}

func makePlatformProjectEntity(r *PlatformProject, m map[string]interface{}, marshal func(interface{}) ([]byte, error), unmarshal func([]byte, any) error) error {
	if value, ok := m["id"].(string); ok {
		r.Id = value
	}

	if value, ok := m["name"].(string); ok {
		r.Name = value
	}

	if value, ok := m["url"].(string); ok {
		r.Url = value
	}

	if value, ok := m["property"].(map[string]interface{}); ok {
		property := make(map[string]string, len(value))
		for key, v := range value {
			if d, ok := v.(string); ok {
				property[key] = d
			}
		}
		r.Property = property
	}

	if value, ok := m["webhooks"].(primitive.A); ok {
		var webhooks []Webhook
		for _, item := range value {
			jsonBytes, err := marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal item: %v", err)
			}

			var webhook Webhook
			if err := unmarshal(jsonBytes, &webhook); err != nil {
				return fmt.Errorf("failed to unmarshal item to Webhook: %v", err)
			}

			webhooks = append(webhooks, webhook)
		}
		r.Webhooks = webhooks
	}

	return nil
}

func makePlatformProjectMap(r *PlatformProject) map[string]interface{} {
	m := map[string]interface{}{
		"id":       r.Id,
		"name":     r.Name,
		"url":      r.Url,
		"property": r.Property,
		"webhooks": r.Webhooks,
	}

	return m
}
