package project

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *Project) MarshalJSON() ([]byte, error) {
	return json.Marshal(makeMap(r))
}

func (r *Project) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return makeEntity(r, m, json.Marshal, json.Unmarshal)
}

func (r *Project) MarshalBSON() ([]byte, error) {
	return bson.Marshal(makeMap(r))
}

func (r *Project) UnmarshalBSON(data []byte) error {
	var m map[string]interface{}
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	return makeEntity(r, m, bson.Marshal, bson.Unmarshal)
}

func makeEntity(r *Project, m map[string]interface{}, marshal func(interface{}) ([]byte, error), unmarshal func([]byte, any) error) error {
	if value, ok := m["id"].(string); ok {
		r.Id = value
	}

	if value, ok := m["name"].(string); ok {
		r.Name = value
	}

	if value, ok := m["description"].(string); ok {
		r.Description = value
	}

	if value, ok := m["start_date"].(string); ok {
		if t, err := time.Parse(time.RFC3339, value); err == nil {
			r.StartDate = &t
		}
	}

	if value, ok := m["end_date"].(string); ok {
		if t, err := time.Parse(time.RFC3339, value); err == nil {
			r.EndDate = &t
		}
	}

	if state, ok := m["state"].(string); ok {
		r.State = GetProjectState(state)
	}

	if value, ok := m["platforms"].(primitive.A); ok {
		var platforms []ProjectPlatform
		for _, item := range value {
			jsonBytes, err := marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal item: %v", err)
			}

			var platform ProjectPlatform
			if err := unmarshal(jsonBytes, &platform); err != nil {
				return fmt.Errorf("failed to unmarshal item to ProjectPlatform: %v", err)
			}

			platforms = append(platforms, platform)
		}
		r.Platforms = platforms
	}

	if value, ok := m["designs"].(primitive.A); ok {
		var designs []ProjectDesign
		for _, item := range value {
			jsonBytes, err := marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal item: %v", err)
			}

			var design ProjectDesign
			if err := unmarshal(jsonBytes, &design); err != nil {
				return fmt.Errorf("failed to unmarshal item to ProjectDesign: %v", err)
			}

			designs = append(designs, design)
		}
		r.Designs = designs
	}

	return nil
}

func makeMap(r *Project) map[string]interface{} {
	m := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"description": r.Description,
		"state":       r.State.String(),
		"start_date":  r.StartDate.Format(time.RFC3339),
		"end_date":    r.EndDate.Format(time.RFC3339),
		"platforms":   r.Platforms,
		"designs":     r.Designs,
	}

	return m
}
