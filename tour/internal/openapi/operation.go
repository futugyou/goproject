package openapi

import "reflect"

type OpenAPIOperation struct {
	Method      string
	Path        string
	Description string
	Request     interface{}
	Response    interface{}
}

func NewOpenAPIOperation(method string, path string, description string, req reflect.Type, resp reflect.Type) *OpenAPIOperation {
	r := &OpenAPIOperation{
		Method:      method,
		Path:        path,
		Description: description,
		Request:     nil,
		Response:    nil,
	}
	if req != nil {
		r.Request = reflect.New(req).Interface()
	}
	if resp != nil {
		r.Response = reflect.New(resp).Interface()
	}
	return r
}
