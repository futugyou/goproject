package pkg

type OpenaiError struct {
	ErrorMessage string `json:"message,omitempty"`
	ErrorType    string `json:"type,omitempty"`
	Param        string `json:"param,omitempty"`
	ErrorCode    string `json:"code,omitempty"`
}
