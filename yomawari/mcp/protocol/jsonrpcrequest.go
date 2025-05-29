package protocol

import "fmt"

type JsonRpcRequest struct {
	JsonRpc          string     `json:"jsonrpc"`
	Method           string     `json:"method"`
	Params           any        `json:"params,omitempty"`
	Id               *RequestId `json:"id"`
	RelatedTransport ITransport `json:"-"`
}

// GetRelatedTransport implements IJsonRpcMessageWithId.
func (j *JsonRpcRequest) GetRelatedTransport() ITransport {
	return j.RelatedTransport
}

// SetRelatedTransport implements IJsonRpcMessageWithId.
func (j *JsonRpcRequest) SetRelatedTransport(transport ITransport) {
	j.RelatedTransport = transport
}

func NewJsonRpcRequest(method string, params any, id *RequestId) *JsonRpcRequest {
	return &JsonRpcRequest{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      id,
	}
}

// GetId implements IJsonRpcMessageWithId.
func (j *JsonRpcRequest) GetId() *RequestId {
	return j.Id
}

// GetJsonRpc implements IJsonRpcMessageWithId.
func (j *JsonRpcRequest) GetJsonRpc() string {
	return "2.0"
}

func (m *JsonRpcRequest) Get(key string) string {
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

func (m *JsonRpcRequest) Set(key, value string) {
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

func (m *JsonRpcRequest) Keys() []string {
	var keys []string
	if mp, ok := m.Params.(map[string]interface{}); ok {
		for k := range mp {
			keys = append(keys, k)
		}
	}
	return keys
}

var _ IJsonRpcMessageWithId = (*JsonRpcRequest)(nil)
