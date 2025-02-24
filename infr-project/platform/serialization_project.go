package platform

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	if value, ok := m["description"].(string); ok {
		r.Description = value
	}

	if value, ok := m["provider_project_id"].(string); ok {
		r.ProviderProjectId = value
	}

	if value, ok := m["image_url"].(string); ok {
		r.ImageUrl = value
	}

	if value, ok := m["image_data"].(primitive.Binary); ok {
		r.ImageData = value.Data
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
	properties := make([]Property, 0, len(r.Properties))
	for _, k := range r.Properties {
		properties = append(properties, k)
	}

	secrets := make([]Secret, 0, len(r.Secrets))
	for _, k := range r.Secrets {
		secrets = append(secrets, k)
	}
	m := map[string]interface{}{
		"id":                  r.Id,
		"name":                r.Name,
		"url":                 r.Url,
		"properties":          properties,
		"secrets":             secrets,
		"webhooks":            r.Webhooks,
		"provider_project_id": r.ProviderProjectId,
		"description":         r.Description,
		"image_data":          r.ImageData,
		"image_url":           r.ImageUrl,
	}

	return m
}
