package errcode

import "fmt"

type Error struct {
	code int
	msg  string
}

var _codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("error code %d already added.", code))
	}
	_codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error code: %d, error msg: %s", e.code, e.msg)
}

func (e *Error) Code() int {
	return e.code
}
func (e *Error) Msg() string {
	return e.msg
}
