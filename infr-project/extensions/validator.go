package extensions

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("lenRange", validateLengthRange)
}

func validateLengthRange(fl validator.FieldLevel) bool {
	param := fl.Param()
	params := strings.Split(param, ",")
	if len(params) != 2 {
		return false
	}

	minLength, err := strconv.Atoi(params[0])
	if err != nil {
		return false
	}
	maxLength, err := strconv.Atoi(params[1])
	if err != nil {
		return false
	}

	name, ok := fl.Field().Interface().(*string)
	if !ok {
		return false
	}
	if name == nil {
		return true
	}
	if len(*name) >= minLength && len(*name) <= maxLength {
		return true
	}
	return false
}
