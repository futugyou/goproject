package lib

import (
	"fmt"
	"strings"
)

type OpenaiError struct {
	ErrorMessage string `json:"message"`
	ErrorType    string `json:"type"`
	Param        string `json:"param"`
	ErrorCode    string `json:"code"`
}

func MessageError(message string) *OpenaiError {
	return &OpenaiError{
		ErrorMessage: message,
		ErrorType:    "invalid parameters",
	}
}

func UnsupportedTypeError(field, value string, list []string) *OpenaiError {
	message := fmt.Sprintf("%s only support [%s], but current value is: %s.", field, strings.Join(list, ","), value)

	return &OpenaiError{
		ErrorMessage: message,
		ErrorType:    "invalid parameters",
		Param:        fmt.Sprintf("current value is: %s", value),
	}
}
