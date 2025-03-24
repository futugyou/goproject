package configuration

import (
	"errors"
	"strings"
	"unicode"
)

const (
	APIKeyAuthType     = "APIKey"
	AccessKeyMinLength = 32
)

var ValidChars = map[rune]struct{}{
	',': {}, '.': {}, ';': {}, ':': {}, '_': {}, '-': {},
	'!': {}, '@': {}, '#': {}, '$': {}, '^': {}, '*': {}, '~': {}, '=': {}, '|': {},
	'[': {}, ']': {}, '{': {}, '}': {}, '(': {}, ')': {},
}

type ServiceAuthorizationConfig struct {
	Enabled            bool   `json:"enabled"`
	AuthenticationType string `json:"authenticationType"`
	HttpHeaderName     string `json:"httpHeaderName"`
	AccessKey1         string `json:"accessKey1"`
	AccessKey2         string `json:"accessKey2"`
}

func NewServiceAuthorizationConfig() *ServiceAuthorizationConfig {
	return &ServiceAuthorizationConfig{
		Enabled:            false,
		AuthenticationType: APIKeyAuthType,
		HttpHeaderName:     "Authorization",
		AccessKey1:         "",
		AccessKey2:         "",
	}
}

func (c *ServiceAuthorizationConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.AuthenticationType != APIKeyAuthType {
		return errors.New("authorization type is not supported, please use 'APIKey'")
	}

	if strings.TrimSpace(c.HttpHeaderName) == "" {
		return errors.New("the HTTP header name cannot be empty")
	}

	if err := validateAccessKey(c.AccessKey1, 1); err != nil {
		return err
	}

	if err := validateAccessKey(c.AccessKey2, 2); err != nil {
		return err
	}

	return nil
}

func validateAccessKey(key string, keyNumber int) error {
	if len(key) == 0 {
		return errors.New("Access Key " + string(rune(keyNumber)) + " is empty.")
	}

	if len(key) < AccessKeyMinLength {
		return errors.New("Access Key " + string(rune(keyNumber)) + " is too short, use at least " + string(rune(AccessKeyMinLength)) + " chars.")
	}

	for _, c := range key {
		if !isValidChar(c) {
			return errors.New("Access Key " + string(rune(keyNumber)) + " contains some invalid chars")
		}
	}

	return nil
}

func isValidChar(c rune) bool {
	if unicode.IsLetter(c) || unicode.IsDigit(c) {
		return true
	}
	_, exists := ValidChars[c]
	return exists
}
