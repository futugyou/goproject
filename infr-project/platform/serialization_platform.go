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
		properties := make(map[string]PropertyInfo)
		for _, item := range value {
			jsonBytes, err := marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal PropertyInfo item: %v", err)
			}

			var proinfo PropertyInfo
			if err := unmarshal(jsonBytes, &proinfo); err != nil {
				return fmt.Errorf("failed to unmarshal item to PropertyInfo: %v", err)
			}

			properties[proinfo.Key] = proinfo
		}
		r.Properties = properties
	}

	if value, ok := m["secrets"].(primitive.A); ok {
		secrets := make(map[string]Secret)
		for _, item := range value {
			jsonBytes, err := marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal Secret item: %v", err)
			}

			var secret Secret
			if err := unmarshal(jsonBytes, &secret); err != nil {
				return fmt.Errorf("failed to unmarshal item to Secret: %v", err)
			}

			secrets[secret.Key] = secret
		}
		r.Secrets = secrets
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

	properties := make([]PropertyInfo, 0, len(r.Properties))
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
