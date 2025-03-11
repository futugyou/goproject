package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

type FunctionDescriptor struct {
	Func             reflect.Value
	Name             string
	Description      string
	JSONSchema       map[string]interface{}
	ParamMarshallers []func(args map[string]interface{}, ctx context.Context) (reflect.Value, error)
	ReturnMarshaller func(results []reflect.Value, ctx context.Context) (interface{}, error)
}

// descriptorCache
var (
	descriptorCache sync.Map // key: string, value: *FunctionDescriptor
)

// GetOrCreateDescriptor
func GetOrCreateDescriptor(fn interface{}, options AIFunctionFactoryOptions) (*FunctionDescriptor, error) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return nil, errors.New("provided fn is not a function")
	}

	// use func name as a part of cache key
	fnPtr := runtime.FuncForPC(v.Pointer()).Name()
	key := fnPtr
	if options.Name != "" {
		key += "_" + options.Name
	}
	if options.Description != "" {
		key += "_" + options.Description
	}

	if cached, ok := descriptorCache.Load(key); ok {
		return cached.(*FunctionDescriptor), nil
	}

	fd, err := newFunctionDescriptor(v, options)
	if err != nil {
		return nil, err
	}
	descriptorCache.Store(key, fd)
	return fd, nil
}

// newFunctionDescriptor
func newFunctionDescriptor(fn reflect.Value, options AIFunctionFactoryOptions) (*FunctionDescriptor, error) {
	fnType := fn.Type()
	numIn := fnType.NumIn()
	paramMarshallers := make([]func(args map[string]interface{}, ctx context.Context) (reflect.Value, error), numIn)

	// prepare parameter name. if it not enough, use param%d
	paramNames := options.ParameterNames
	if len(paramNames) < numIn {
		for i := len(paramNames); i < numIn; i++ {
			paramNames = append(paramNames, fmt.Sprintf("param%d", i))
		}
	}

	// prepare convertion for each parameter
	for i := 0; i < numIn; i++ {
		inType := fnType.In(i)
		paramName := paramNames[i]
		// If the parameter implements context.Context, the passed context is used directly
		if inType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			paramMarshallers[i] = func(args map[string]interface{}, ctx context.Context) (reflect.Value, error) {
				return reflect.ValueOf(ctx), nil
			}
		} else {
			// Otherwise, find the corresponding value from the parameter map and perform necessary type conversion
			// closure is used to bind paramName and the target type
			paramMarshallers[i] = func(paramName string, targetType reflect.Type) func(args map[string]interface{}, ctx context.Context) (reflect.Value, error) {
				return func(args map[string]interface{}, ctx context.Context) (reflect.Value, error) {
					val, ok := args[paramName]
					if !ok {
						// If no parameter is provided, the zero value for that type is returned.
						return reflect.Zero(targetType), nil
					}
					rv := reflect.ValueOf(val)
					if rv.Type().AssignableTo(targetType) {
						return rv, nil
					}
					// Try converting types via JSON encoding/decoding
					data, err := json.Marshal(val)
					if err != nil {
						return reflect.Zero(targetType), fmt.Errorf("failed to marshal parameter %s: %w", paramName, err)
					}
					ptr := reflect.New(targetType)
					err = json.Unmarshal(data, ptr.Interface())
					if err != nil {
						return reflect.Zero(targetType), fmt.Errorf("failed to unmarshal parameter %s: %w", paramName, err)
					}
					return ptr.Elem(), nil
				}
			}(paramName, inType)
		}
	}

	// Constructing a return value converter
	var returnMarshaller func(results []reflect.Value, ctx context.Context) (interface{}, error)
	numOut := fnType.NumOut()
	if numOut == 0 {
		returnMarshaller = func(results []reflect.Value, ctx context.Context) (interface{}, error) {
			return nil, nil
		}
	} else if numOut == 1 {
		returnMarshaller = func(results []reflect.Value, ctx context.Context) (interface{}, error) {
			return results[0].Interface(), nil
		}
	} else if numOut == 2 {
		// Assume the second return value is error
		returnMarshaller = func(results []reflect.Value, ctx context.Context) (interface{}, error) {
			errVal := results[1]
			if !errVal.IsNil() {
				errInterface := errVal.Interface()
				if err, ok := errInterface.(error); ok {
					return nil, err
				}
			}
			return results[0].Interface(), nil
		}
	} else {
		return nil, errors.New("function has unsupported number of return values")
	}

	// Determine the function name: if not provided by the user, the name obtained by reflection is used
	name := options.Name
	if name == "" {
		name = runtime.FuncForPC(fn.Pointer()).Name()
	}

	description := options.Description

	// Generate a simple JSON Schema description (here just a simple description of the type of each parameter)
	jsonSchema := generateJSONSchema(fnType, name, description, paramNames)

	return &FunctionDescriptor{
		Func:             fn,
		Name:             name,
		Description:      description,
		JSONSchema:       jsonSchema,
		ParamMarshallers: paramMarshallers,
		ReturnMarshaller: returnMarshaller,
	}, nil
}

// generateJSONSchema simply generates a JSON Schema to describe the input parameters of the function
func generateJSONSchema(fnType reflect.Type, name, description string, paramNames []string) map[string]interface{} {
	properties := make(map[string]interface{})
	numIn := fnType.NumIn()
	for i := 0; i < numIn; i++ {
		inType := fnType.In(i)
		// Parameters that implement context.Context are not reflected in the Schema
		if inType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			continue
		}
		paramName := paramNames[i]
		properties[paramName] = map[string]string{
			"type": goTypeToJSONType(inType),
		}
	}
	schema := map[string]interface{}{
		"title":       name,
		"description": description,
		"type":        "object",
		"properties":  properties,
	}
	return schema
}

// goTypeToJSONType maps Go types to type descriptions in JSON Schema
func goTypeToJSONType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	default:
		return "string"
	}
}

// Invoke calls the function based on the passed parameters map and context and returns the converted result
func (fd *FunctionDescriptor) Invoke(args map[string]interface{}, ctx context.Context) (interface{}, error) {
	fnType := fd.Func.Type()
	numIn := fnType.NumIn()
	inValues := make([]reflect.Value, numIn)
	for i, marshaller := range fd.ParamMarshallers {
		val, err := marshaller(args, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameter %d: %w", i, err)
		}
		inValues[i] = val
	}
	results := fd.Func.Call(inValues)
	return fd.ReturnMarshaller(results, ctx)
}
