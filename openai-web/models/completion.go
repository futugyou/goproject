package models

import (
	"github.com/beego/beego/v2/core/validation"
)

type CompletionModel struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	MaxTokens        int64    `json:"max_tokens"`
	Temperature      float64  `json:"temperature" valid:"ValidateTemperature"`
	Top_P            float64  `json:"top_p" valid:"ValidateTemperature"`
	FrequencyPenalty float64  `json:"frequency_penalty" valid:"ValidatePenalty"`
	PresencePenalty  float64  `json:"presence_penalty" valid:"ValidatePenalty"`
	BestOf           int64    `json:"best_of" valid:"Range(1, 20)"`
	Echo             bool     `json:"echo"`
	Logprobs         int      `json:"logprobs"`
	Stop             []string `json:"stop"`
	Stream           bool     `json:"stream"`
	Suffix           string   `json:"suffix,omitempty"`
}

func init() {
	validation.AddCustomFunc("ValidateTemperature", validateTemperature)
	validation.AddCustomFunc("ValidatePenalty", validatePenalty)
}

var validateTemperature validation.CustomFunc = func(v *validation.Validation, obj interface{}, key string) {
	temperature, ok := obj.(float64)

	if !ok {
		return
	}

	if temperature < 0.0 || temperature > 1.0 {
		v.AddError(key, "must in 0~1")
	}
}

var validatePenalty validation.CustomFunc = func(v *validation.Validation, obj interface{}, key string) {
	temperature, ok := obj.(float64)

	if !ok {
		return
	}

	if temperature < 0.0 || temperature > 2.0 {
		v.AddError(key, "must in 0~2")
	}
}
