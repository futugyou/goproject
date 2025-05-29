package protocol

import (
	"encoding/json"
	"fmt"
)

type JsonRpcNotification struct {
	JsonRpc          string          `json:"jsonrpc"`
	Method           string          `json:"method"`
	Params           json.RawMessage `json:"params,omitempty"`
	RelatedTransport ITransport      `json:"-"`
}

// GetRelatedTransport implements IJsonRpcMessage.
func (m *JsonRpcNotification) GetRelatedTransport() ITransport {
	return m.RelatedTransport
}

// SetRelatedTransport implements IJsonRpcMessage.
func (m *JsonRpcNotification) SetRelatedTransport(transport ITransport) {
	m.RelatedTransport = transport
}

func NewJsonRpcNotification(method string, params json.RawMessage) *JsonRpcNotification {
	return &JsonRpcNotification{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

func NewJsonRpcNotificationWithTransport(method string, params json.RawMessage, transport ITransport) *JsonRpcNotification {
	return &JsonRpcNotification{
		JsonRpc:          "2.0",
		Method:           method,
		Params:           params,
		RelatedTransport: transport,
	}
}

func (m *JsonRpcNotification) Get(key string) string {
	if m.Params == nil {
		return ""
	}
	var mp map[string]interface{}
	if err := json.Unmarshal(m.Params, &mp); err != nil {
		return ""
	}
	if val, exists := mp[key]; exists {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func (m *JsonRpcNotification) Set(key, value string) {
	if m.Params == nil {
		m.Params = json.RawMessage{}
	}
	var mp map[string]interface{}
	if err := json.Unmarshal(m.Params, &mp); err != nil {
		return
	}
	mp[key] = value
	d, err := json.Marshal(mp)
	if err != nil {
		return
	}
	m.Params = d
}

func (m *JsonRpcNotification) Keys() []string {
	var keys []string
	if m.Params == nil {
		return keys
	}
	var mp map[string]interface{}
	if err := json.Unmarshal(m.Params, &mp); err != nil {
		return keys
	}
	for k := range mp {
		keys = append(keys, k)
	}
	return keys
}

// GetJsonRpc implements IJsonRpcMessage.
func (j *JsonRpcNotification) GetJsonRpc() string {
	return "2.0"
}

var _ IJsonRpcMessage = (*JsonRpcNotification)(nil)
