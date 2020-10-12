package errcode

import (
	"fmt"
	"net/http"
)

var (
	Success                   = NewError(0, "success")
	ServiceError              = NewError(1000000, "ServiceError")
	InvalidParams             = NewError(1000001, "InvalidParams")
	NotFound                  = NewError(1000002, "NotFound")
	UnauthorizedAuthNotExist  = NewError(1000003, "UnauthorizedAuthNotExist")
	UnauthorizedTokenError    = NewError(1000004, "UnauthorizedTokenError")
	UnauthorizedTokenTimeout  = NewError(1000005, "UnauthorizedTokenTimeout")
	UnauthorizedTokenGenerate = NewError(1000006, "UnauthorizedTokenGenerate")
	TooManyRequests           = NewError(1000007, "TooManyRequests")

	ErrorGetListFail   = NewError(2000003, "ErrorGetListFail")
	ErrorCreateTagFail = NewError(2000004, "ErrorCreateTagFail")
	ErrorUpdateTagFail = NewError(2000005, "ErrorUpdateTagFail")
	ErrorDeleteTagFail = NewError(2000006, "ErrorDeleteTagFail")
	ErrorCountTagFail  = NewError(2000007, "ErrorCountTagFail")

	ErrorGetArticleFail    = NewError(2000012, "ErrorGetArticleFail")
	ErrorGetArticlesFail   = NewError(2000013, "ErrorGetArticlesFail")
	ErrorCreateArticleFail = NewError(2000014, "ErrorCreateArticleFail")
	ErrorUpdateArticleFail = NewError(2000015, "ErrorUpdateArticleFail")
	ErrorDeleteArticleFail = NewError(2000016, "ErrorDeleteArticleFail")

	ErrorUploadFileFail = NewError(2000026, "ErrorUploadFileFail")
)

type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("readly have thi code : %d", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code : %d,  message: %s,", e.Code(), e.Msg())
}
func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args)
}

func (e *Error) Details() []string {
	return e.details
}
func (e *Error) WithDetails(details ...string) *Error {
	e.details = []string{}
	for _, d := range details {
		e.details = append(e.details, d)
	}
	return e
}

func (e *Error) StatusCode() int {
	switch e.code {
	case Success.code:
		return http.StatusOK
	case ServiceError.code:
		return http.StatusInternalServerError
	case InvalidParams.code:
		return http.StatusBadRequest
	case NotFound.code:
		return http.StatusNotFound
	case UnauthorizedAuthNotExist.code:
		fallthrough
	case UnauthorizedTokenError.code:
		fallthrough
	case UnauthorizedTokenTimeout.code:
		fallthrough
	case UnauthorizedTokenGenerate.code:
		return http.StatusUnauthorized
	case TooManyRequests.code:
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
