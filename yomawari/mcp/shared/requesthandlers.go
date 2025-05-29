package shared

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type RequestHandler func(ctx context.Context, request *protocol.JsonRpcRequest) (json.RawMessage, error)
type GenericRequestHandler[TRequest any, TResponse any] func(ctx context.Context, request *TRequest, tran protocol.ITransport) (*TResponse, error)
type RequestUnmarshaler[TRequest any] func(data interface{}) (*TRequest, error)
type RepsonseMarshaler[TResponse any] func(data TResponse) (json.RawMessage, error)

func DefaultJsonUnmarshaler[TRequest any](data interface{}) (*TRequest, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var req TRequest
	err = json.Unmarshal(bytes, &req)
	return &req, err
}

func DefaultRepsonseMarshaler[TResponse any](data TResponse) (json.RawMessage, error) {
	return json.Marshal(data)
}

type RequestHandlers struct {
	handlers map[string]RequestHandler
	mu       sync.RWMutex
}

func NewRequestHandlers() *RequestHandlers {
	return &RequestHandlers{
		handlers: make(map[string]RequestHandler),
	}
}

func (c *RequestHandlers) Count() int {
	return len(c.handlers)
}

func (c *RequestHandlers) IsEmpty() bool {
	if c == nil || len(c.handlers) == 0 {
		return true
	}
	return false
}

func (c *RequestHandlers) Get(method string) (RequestHandler, bool) {
	v, ok := c.handlers[method]
	return v, ok
}

func (c *RequestHandlers) Clear() {
	c.handlers = make(map[string]RequestHandler)
}

func (c *RequestHandlers) Add(method string, handler RequestHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.handlers[method]; ok {
		return
	}
	c.handlers[method] = handler
}

func (c *RequestHandlers) TryAdd(method string, handler RequestHandler) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.handlers[method]; ok {
		return false
	}
	c.handlers[method] = handler
	return true
}

func (c *RequestHandlers) Remove(method string) {
	delete(c.handlers, method)
}

func (c *RequestHandlers) Contains(method string) bool {
	if _, ok := c.handlers[method]; ok {
		return true
	}
	return false
}

func GenericRequestHandlerAdd[TRequest any, TResponse any](
	handers *RequestHandlers,
	method string,
	handler GenericRequestHandler[TRequest, TResponse],
	unmarshaler RequestUnmarshaler[TRequest],
	marshaler RepsonseMarshaler[TResponse],
) {
	if handers == nil {
		return
	}

	if unmarshaler == nil {
		unmarshaler = DefaultJsonUnmarshaler[TRequest]
	}
	if marshaler == nil {
		marshaler = DefaultRepsonseMarshaler[TResponse]
	}

	handers.mu.Lock()
	defer handers.mu.Unlock()
	handers.handlers[method] = func(ctx context.Context, request *protocol.JsonRpcRequest) (json.RawMessage, error) {
		requestBody := request.Params
		req, err := unmarshaler(requestBody)
		if err != nil {
			return nil, err
		}
		resp, err := handler(ctx, req, request.RelatedTransport)
		if err != nil {
			return nil, err
		}
		return marshaler(*resp)
	}
}
