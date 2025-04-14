package messages

import "fmt"

type JsonRpcNotification struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

func NewJsonRpcNotification(method string, params any) *JsonRpcNotification {
	return &JsonRpcNotification{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

func (m *JsonRpcNotification) Get(key string) string {
	if m.Params == nil {
		return ""
	}

	if mp, ok := m.Params.(map[string]interface{}); ok {
		if val, exists := mp[key]; exists {
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

func (m *JsonRpcNotification) Set(key, value string) {
	var mp map[string]interface{}

	if m.Params != nil {
		if existingMap, ok := m.Params.(map[string]interface{}); ok {
			mp = existingMap
		}
	}

	if mp == nil {
		mp = make(map[string]interface{})
	}

	mp[key] = value
	m.Params = mp
}

func (m *JsonRpcNotification) Keys() []string {
	var keys []string
	if mp, ok := m.Params.(map[string]interface{}); ok {
		for k := range mp {
			keys = append(keys, k)
		}
	}
	return keys
}

// GetJsonRpc implements IJsonRpcMessage.
func (j *JsonRpcNotification) GetJsonRpc() string {
	return "2.0"
}

var _ IJsonRpcMessage = (*JsonRpcNotification)(nil)
