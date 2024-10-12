package platform

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	if value, ok := m["properties"].(primitive.A); ok {
		properties, err := parseArrayToMap[Property](value, marshal, unmarshal)
		if err != nil {
			return err
		}
		r.Properties = properties
	}

	if value, ok := m["secrets"].(primitive.A); ok {
		secrets, err := parseArrayToMap[Secret](value, marshal, unmarshal)
		if err != nil {
			return err
		}
		r.Secrets = secrets
	}

	if value, ok := m["projects"].(primitive.A); ok {
		projects, err := parseArrayToMap[PlatformProject](value, marshal, unmarshal)
		if err != nil {
			return err
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

	properties := make([]Property, 0, len(r.Properties))
	for _, k := range r.Properties {
		properties = append(properties, k)
	}

	secrets := make([]Secret, 0, len(r.Secrets))
	for _, k := range r.Secrets {
		secrets = append(secrets, k)
	}

	m := map[string]interface{}{
		"id":         r.Id,
		"name":       r.Name,
		"activate":   r.Activate,
		"url":        r.Url,
		"properties": properties,
		"secrets":    secrets,
		"projects":   projects,
		"tags":       r.Tags,
		"is_deleted": r.IsDeleted,
	}

	return m
}

func parseArrayToMap[T any](array []any, marshalFunc func(any) ([]byte, error), unmarshalFunc func([]byte, any) error) (map[string]T, error) {
	result := make(map[string]T)
	for _, item := range array {
		jsonBytes, err := marshalFunc(item)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal item: %v", err)
		}

		var obj T
		if err := unmarshalFunc(jsonBytes, &obj); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %v", err)
		}

		keyGetter, ok := any(&obj).(interface{ GetKey() string })
		if !ok {
			return nil, fmt.Errorf("type %T does not implement GetKey method", obj)
		}

		result[keyGetter.GetKey()] = obj
	}
	return result, nil
}
