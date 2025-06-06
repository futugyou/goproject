package resource

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MarshalJSON is a custom marshaler for Resource that handles the serialization of ResourceType.
// In this case, we can skip MarshalJSON, only implement UnmarshalJSON
func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(makeMap(r))
}

// UnmarshalJSON is a custom unmarshaler for Resource that handles the deserialization of ResourceType.
func (r *Resource) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return makeEntity(r, m)
}

func (r Resource) MarshalBSON() ([]byte, error) {
	return bson.Marshal(makeMap(r))

	// type Alias Resource
	// return json.Marshal(&struct {
	// 	Type string `json:"type"`
	// 	*Alias
	// }{
	// 	Type:  r.Type.String(),
	// 	Alias: (*Alias)(r),
	// })
}

func (r *Resource) UnmarshalBSON(data []byte) error {
	var m map[string]interface{}
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	return makeEntity(r, m)

	// type Alias Resource
	// aux := &struct {
	// 	Type string `json:"type"`
	// 	*Alias
	// }{
	// 	Alias: (*Alias)(r),
	// }

	// if err := json.Unmarshal(data, &aux); err != nil {
	// 	return err
	// }

	// switch aux.Type {
	// case string(DrawIO):
	// 	r.Type = DrawIO
	// case string(Markdown):
	// 	r.Type = Markdown
	// case string(Excalidraw):
	// 	r.Type = Excalidraw
	// case string(Plate):
	// 	r.Type = Plate
	// default:
	// 	return json.Unmarshal(data, &r)
	// }

	// return nil
}

func makeEntity(r *Resource, m map[string]interface{}) error {
	if id, ok := m["id"].(string); ok {
		r.Id = id
	}

	if name, ok := m["name"].(string); ok {
		r.Name = name
	}

	if value, ok := m["is_deleted"].(bool); ok {
		r.IsDeleted = value
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

	if v, ok := m["version"]; ok {
		switch version := v.(type) {
		case int:
			r.Version = version
		case int32:
			r.Version = int(version)
		default:
			r.Version = 1
		}
	}

	resourceType, _ := m["type"].(string)
	r.Type = GetResourceType(resourceType)

	if data, ok := m["data"].(string); ok {
		r.Data = data
	}

	if data, ok := m["imageData"].(string); ok {
		r.ImageData = data
	}

	if createdAt, ok := m["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			r.CreatedAt = t
		} else {
			return err
		}
	}

	if createdAt, ok := m["updated_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			r.UpdatedAt = t
		} else {
			return err
		}
	}

	return nil
}

func makeMap(r Resource) map[string]interface{} {
	m := map[string]interface{}{
		"id":         r.Id,
		"name":       r.Name,
		"version":    r.Version,
		"data":       r.Data,
		"imageData":  r.ImageData,
		"tags":       r.Tags,
		"created_at": r.CreatedAt.Format(time.RFC3339),
		"updated_at": r.UpdatedAt.Format(time.RFC3339),
		"is_deleted": r.IsDeleted,
	}
	if r.Type != nil {
		m["type"] = r.Type.String()
	} else {
		m["type"] = string(Markdown)
	}

	return m
}
