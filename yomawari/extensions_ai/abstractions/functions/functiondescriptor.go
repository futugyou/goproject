package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/utilities"
)

type FunctionDescriptor struct {
	Func             reflect.Value
	Name             string
	Description      string
	JSONSchema       map[string]interface{}
	ParamMarshallers []func(ctx context.Context, args AIFunctionArguments) (reflect.Value, error)
	ReturnMarshaller func(ctx context.Context, results []reflect.Value) (interface{}, error)
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
	paramMarshallers := make([]func(ctx context.Context, args AIFunctionArguments) (reflect.Value, error), numIn)

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
			paramMarshallers[i] = func(ctx context.Context, args AIFunctionArguments) (reflect.Value, error) {
				return reflect.ValueOf(ctx), nil
			}
		} else {
			// Otherwise, find the corresponding value from the parameter map and perform necessary type conversion
			// closure is used to bind paramName and the target type
			paramMarshallers[i] = func(paramName string, targetType reflect.Type) func(ctx context.Context, args AIFunctionArguments) (reflect.Value, error) {
				return func(ctx context.Context, args AIFunctionArguments) (reflect.Value, error) {
					val, ok := args.Get(paramName)
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
	var returnMarshaller func(ctx context.Context, results []reflect.Value) (interface{}, error)
	numOut := fnType.NumOut()
	if numOut == 0 {
		returnMarshaller = func(ctx context.Context, results []reflect.Value) (interface{}, error) {
			return nil, nil
		}
	} else if numOut == 1 {
		returnMarshaller = func(ctx context.Context, results []reflect.Value) (interface{}, error) {
			return results[0].Interface(), nil
		}
	} else if numOut == 2 {
		// Assume the second return value is error
		returnMarshaller = func(ctx context.Context, results []reflect.Value) (interface{}, error) {
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
	jsonSchema, err := utilities.CreateFunctionJsonSchema(fnType, name, description, paramNames, utilities.DefaultAIJsonSchemaCreateOptions)
	if err != nil {
		return nil, err
	}

	return &FunctionDescriptor{
		Func:             fn,
		Name:             name,
		Description:      description,
		JSONSchema:       jsonSchema,
		ParamMarshallers: paramMarshallers,
		ReturnMarshaller: returnMarshaller,
	}, nil
}

// Invoke calls the function based on the passed parameters map and context and returns the converted result
func (fd *FunctionDescriptor) Invoke(ctx context.Context, args AIFunctionArguments) (interface{}, error) {
	fnType := fd.Func.Type()
	numIn := fnType.NumIn()
	inValues := make([]reflect.Value, numIn)
	for i, marshaller := range fd.ParamMarshallers {
		val, err := marshaller(ctx, args)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameter %d: %w", i, err)
		}
		inValues[i] = val
	}
	results := fd.Func.Call(inValues)
	return fd.ReturnMarshaller(ctx, results)
}
