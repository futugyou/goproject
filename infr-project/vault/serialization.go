package vault

import (
	"encoding/json"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	tool "github.com/futugyou/extensions"
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
		"type_identity": r.TypeIdentity,
		"tags":          r.Tags,
		"extension":     r.Extension,
		"description":   r.Description,
	}

	if value, err := tool.AesCTREncrypt(r.Value, os.Getenv("Encrypt_Key")); err != nil {
		return nil, err
	} else {
		m["value"] = value
	}

	if r.StorageMedia != nil {
		m["storage_media"] = r.StorageMedia.String()
	} else {
		m["storage_media"] = string(StorageMediaLocal)
	}
	if r.VaultType != nil {
		m["vault_type"] = r.VaultType.String()
	} else {
		m["vault_type"] = string(VaultTypeCommon)
	}
	if r.State != nil {
		m["state"] = r.State.String()
	} else {
		m["state"] = string(VaultStateDefault)
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

	if value, ok := m["description"].(string); ok {
		w.Description = value
	}

	if value, ok := m["extension"].(map[string]string); ok {
		w.Extension = value
	}

	if value, ok := m["value"].(string); ok {
		if value, err := tool.AesCTRDecrypt(value, os.Getenv("Encrypt_Key")); err != nil {
			return err
		} else {
			w.Value = value
		}
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

	value, _ := m["storage_media"].(string)
	w.StorageMedia = GetStorageMedia(value)

	value, _ = m["state"].(string)
	w.State = GetVaultState(value)

	value, _ = m["vault_type"].(string)
	w.VaultType = GetVaultType(value)

	return nil
}
