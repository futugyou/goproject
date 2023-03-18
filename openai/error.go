package openai

import (
	"encoding/json"
	"fmt"
	"strings"
)

type OpenaiError struct {
	ErrorMessage string `json:"message"`
	ErrorType    string `json:"type"`
	Param        string `json:"param"`
	ErrorCode    string `json:"code"`
}

func (e *OpenaiError) Error() string {
	if e == nil {
		return ""
	}

	result, _ := json.Marshal(e)
	return string(result)
}

func MessageError(message string) *OpenaiError {
	return &OpenaiError{
		ErrorMessage: message,
		ErrorType:    "invalid parameters",
	}
}

func SystemError(message string) *OpenaiError {
	return &OpenaiError{
		ErrorMessage: message,
		ErrorType:    "system error",
	}
}

func UnsupportedTypeError[T any](field string, value T, list []T) *OpenaiError {
	message := fmt.Sprintf("%s only support [%s], but current value is: %v.", field, enumjoin(list, ","), value)

	return &OpenaiError{
		ErrorMessage: message,
		ErrorType:    "invalid parameters",
		Param:        fmt.Sprintf("current value is: %v", value),
	}
}

func enumjoin[T any](elems []T, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%v", elems[0])
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(fmt.Sprintf("%v", elems[i]))
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(fmt.Sprintf("%v", elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(fmt.Sprintf("%v", s))
	}
	return b.String()
}
