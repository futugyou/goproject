package domains

import (
	"encoding/json"
	"maps"
)

type ApplicationProviderMetadata struct {
	BaseProviderMetadata
	FastfedHandshakeRegisterUri string         `json:"fastfed_handshake_register_uri"`
	OtherParameters             map[string]any `json:"-"`
}

func (a *ApplicationProviderMetadata) UnmarshalJSON(data []byte) error {
	type Alias ApplicationProviderMetadata
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var rawMap map[string]any
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	delete(rawMap, "fastfed_handshake_register_uri")

	a.OtherParameters = rawMap
	return nil
}

func (a ApplicationProviderMetadata) MarshalJSON() ([]byte, error) {
	result := make(map[string]any)

	result["fastfed_handshake_register_uri"] = a.FastfedHandshakeRegisterUri

	maps.Copy(result, a.OtherParameters)

	return json.Marshal(result)
}
