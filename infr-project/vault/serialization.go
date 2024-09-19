package vault

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (w Vault) MarshalJSON() ([]byte, error) {
	return w.commonMarshal(json.Marshal)
}

func (w Vault) MarshalBSON() ([]byte, error) {
	return w.commonMarshal(bson.Marshal)
}

func (r Vault) commonMarshal(marshal func(interface{}) ([]byte, error)) ([]byte, error) {
	m := map[string]interface{}{
		"id":            r.Id,
		"key":           r.Key,
		"value":         r.Value,
		"type_identity": r.TypeIdentity,
		"tags":          r.Tags,
	}
	if r.StorageMedia != nil {
		m["storage_media"] = r.StorageMedia.String()
	}
	if r.VaultType != nil {
		m["vault_type"] = r.VaultType.String()
	}
	return marshal(m)
}

func (w *Vault) UnmarshalJSON(data []byte) error {
	return w.commonUnmarshal(data, json.Unmarshal)
}

func (w *Vault) UnmarshalBSON(data []byte) error {
	return w.commonUnmarshal(data, bson.Unmarshal)
}

func (w *Vault) commonUnmarshal(data []byte, unmarshal func([]byte, any) error) error {
	var m map[string]interface{}
	if err := unmarshal(data, &m); err != nil {
		return err
	}

	if value, ok := m["id"].(string); ok {
		w.Id = value
	}

	if value, ok := m["key"].(string); ok {
		w.Key = value
	}

	if value, ok := m["value"].(string); ok {
		w.Value = value
	}

	if value, ok := m["type_identity"].(string); ok {
		w.TypeIdentity = value
	}

	if value, ok := m["tags"].([]string); ok {
		w.Tags = value
	} else if value, ok := m["tags"].(primitive.A); ok {
		tags := make([]string, 0)
		for _, item := range value {
			if tag, ok := item.(string); ok {
				tags = append(tags, tag)
			}
		}
		w.Tags = tags
	}

	if value, ok := m["storage_media"].(string); ok {
		w.StorageMedia = GetStorageMedia(value)
	}

	if value, ok := m["vault_type"].(string); ok {
		w.VaultType = GetVaultType(value)
	}

	return nil
}
