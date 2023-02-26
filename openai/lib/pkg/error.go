package pkg

type OpenaiError struct {
	ErrorMessage string `json:"message"`
	ErrorType    string `json:"type"`
	Param        string `json:"param"`
	ErrorCode    string `json:"code"`
}
