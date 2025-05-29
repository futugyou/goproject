package protocol

import (
	"encoding/json"
	"errors"
)

type IJsonRpcMessage interface {
	GetJsonRpc() string
	GetRelatedTransport() ITransport
	SetRelatedTransport(transport ITransport)
}

type IJsonRpcMessageWithId interface {
	IJsonRpcMessage
	GetId() *RequestId
}

func UnmarshalJsonRpcMessage(data []byte) (IJsonRpcMessage, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var version string
	if err := json.Unmarshal(raw["jsonrpc"], &version); err != nil || version != "2.0" {
		return nil, errors.New("invalid or missing jsonrpc version")
	}

	_, hasID := raw["id"]
	_, hasMethod := raw["method"]
	_, hasError := raw["error"]

	if hasID && !hasMethod {
		if hasError {
			var msg JsonRpcError
			if err := json.Unmarshal(data, &msg); err != nil {
				return nil, err
			}
			return &msg, nil
		}
		if _, hasResult := raw["result"]; hasResult {
			var msg JsonRpcResponse
			if err := json.Unmarshal(data, &msg); err != nil {
				return nil, err
			}
			return &msg, nil
		}
		return nil, errors.New("response must have either result or error")
	}

	if hasMethod && !hasID {
		var msg JsonRpcNotification
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil
	}

	if hasMethod && hasID {
		var msg JsonRpcRequest
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil
	}

	return nil, errors.New("invalid JSON-RPC message format")
}

func MarshalJsonRpcMessage(msg IJsonRpcMessage) ([]byte, error) {
	switch v := msg.(type) {
	case *JsonRpcRequest:
		return json.Marshal(v)
	case *JsonRpcNotification:
		return json.Marshal(v)
	case *JsonRpcResponse:
		return json.Marshal(v)
	case *JsonRpcError:
		return json.Marshal(v)
	default:
		return nil, errors.New("unknown JSON-RPC message type")
	}
}
