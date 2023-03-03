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

func NewError(value string, list []string) *OpenaiError {
	message := value
	if list != nil {
		message = strings.Join(list, ",")
	}

	return &OpenaiError{
		ErrorMessage: fmt.Sprintf("only support  %s", message),
		ErrorType:    "invalid parameters",
		Param:        fmt.Sprintf("current value is: %s", value),
	}
}
