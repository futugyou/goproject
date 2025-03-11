package functions

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"regexp"
	"sync"

	"github.com/futugyou/ai-extension/abstractions/functions"
)

// AIFunctionFactory provides methods to create AIFunction instances.
type AIFunctionFactory struct {
	defaultOptions AIFunctionFactoryOptions
}

// NewAIFunctionFactory creates a new AIFunctionFactory.
func NewAIFunctionFactory() *AIFunctionFactory {
	return &AIFunctionFactory{
		defaultOptions: AIFunctionFactoryOptions{},
	}
}

// SanitizeMemberName removes characters from a .NET member name that shouldn't be used in an AI function name.
func SanitizeMemberName(memberName string) string {
	if memberName == "" {
		panic("memberName cannot be nil")
	}
	invalidNameCharsRegex := regexp.MustCompile(`[^0-9A-Za-z_]`)
	return invalidNameCharsRegex.ReplaceAllString(memberName, "_")
}

// Create creates an AIFunction instance for a method, specified via a delegate.
func (factory *AIFunctionFactory) Create(method interface{}, options *AIFunctionFactoryOptions) (functions.AIFunction, error) {
	if method == nil {
		return nil, errors.New("method cannot be nil")
	}

	methodValue := reflect.ValueOf(method)
	if methodValue.Kind() != reflect.Func {
		return nil, errors.New("method must be a function")
	}

	opts := factory.defaultOptions
	if options != nil {
		opts = *options
	}

	return &reflectionAIFunction{
		method:  methodValue,
		options: opts,
	}, nil
}

// reflectionAIFunction is an implementation of AIFunction that uses reflection to invoke a method.
type reflectionAIFunction struct {
	functions.BaseAIFunction
	method  reflect.Value
	options AIFunctionFactoryOptions
}

// Invoke invokes the method with the provided arguments.
func (f *reflectionAIFunction) InvokeCore(ctx context.Context, arguments map[string]interface{}) (interface{}, error) {
	methodType := f.method.Type()
	args := make([]reflect.Value, methodType.NumIn())

	for i := 0; i < methodType.NumIn(); i++ {
		argType := methodType.In(i)
		argValue, ok := arguments[argType.Name()]
		if !ok {
			return nil, errors.New("missing argument: " + argType.Name())
		}
		args[i] = reflect.ValueOf(argValue)
	}

	results := f.method.Call(args)
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].Interface(), nil
}

// PooledMemoryStream implements a simple write-only memory stream that uses pooled buffers.
type PooledMemoryStream struct {
	buffer     *bytes.Buffer
	bufferPool sync.Pool
}

func NewPooledMemoryStream(initialCapacity int) *PooledMemoryStream {
	return &PooledMemoryStream{
		buffer: bytes.NewBuffer(make([]byte, 0, initialCapacity)),
		bufferPool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (pms *PooledMemoryStream) Write(data []byte) (int, error) {
	return pms.buffer.Write(data)
}

func (pms *PooledMemoryStream) GetBuffer() []byte {
	return pms.buffer.Bytes()
}

func (pms *PooledMemoryStream) Dispose() {
	pms.buffer.Reset()
	pms.bufferPool.Put(pms.buffer)
	pms.buffer = nil
}
