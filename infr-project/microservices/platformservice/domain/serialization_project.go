package domain

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlatformProject
func (r *PlatformProject) MarshalJSON() ([]byte, error) {
	return json.Marshal(makePlatformProjectMap(r))
}

func (r *PlatformProject) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return makePlatformProjectEntity(r, m, json.Marshal, json.Unmarshal)
}

func (r *PlatformProject) MarshalBSON() ([]byte, error) {
	return bson.Marshal(makePlatformProjectMap(r))
}

func (r *PlatformProject) UnmarshalBSON(data []byte) error {
	var m map[string]any
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	return makePlatformProjectEntity(r, m, bson.Marshal, bson.Unmarshal)
}

func makePlatformProjectEntity(r *PlatformProject, m map[string]any, marshal func(any) ([]byte, error), unmarshal func([]byte, any) error) error {
	if value, ok := m["id"].(string); ok {
		r.ID = value
	}

	if value, ok := m["name"].(string); ok {
		r.Name = value
	}

	if value, ok := m["url"].(string); ok {
		r.Url = value
	}

	if value, ok := m["description"].(string); ok {
		r.Description = value
	}

	if value, ok := m["provider_project_id"].(string); ok {
		r.ProviderProjectId = value
	}

	if value, ok := m["image_url"].(string); ok {
		r.ImageUrl = value
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

	if raw, ok := m["webhook"].(primitive.M); ok {
		bytes, _ := bson.Marshal(raw)
		var webhook Webhook
		_ = bson.Unmarshal(bytes, &webhook)
		r.Webhook = &webhook
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

	return nil
}

func makePlatformProjectMap(r *PlatformProject) map[string]any {
	properties := make([]Property, 0, len(r.Properties))
	for _, k := range r.Properties {
		properties = append(properties, k)
	}

	secrets := make([]Secret, 0, len(r.Secrets))
	for _, k := range r.Secrets {
		secrets = append(secrets, k)
	}
	m := map[string]any{
		"id":                  r.ID,
		"name":                r.Name,
		"url":                 r.Url,
		"properties":          properties,
		"secrets":             secrets,
		"provider_project_id": r.ProviderProjectId,
		"description":         r.Description,
		"image_url":           r.ImageUrl,
		"tags":                r.Tags,
	}

	if r.Webhook != nil {
		m["webhook"] = r.Webhook
	}

	return m
}
