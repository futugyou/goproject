package models

import (
	"github.com/beego/beego/v2/adapter/validation"
)

type QuestionAnswer struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	MaxTokens        int32    `json:"max_tokens"`
	Temperature      float32  `json:"temperature" valid:"ValidateTemperature"`
	TopP             float32  `json:"top_p" valid:"ValidateTemperature"`
	FrequencyPenalty float32  `json:"frequency_penalty" valid:"ValidatePenalty"`
	PresencePenalty  float32  `json:"presence_penalty" valid:"ValidatePenalty"`
	BestOf           int32    `json:"best_of" valid:"Range(1, 20)"`
	Echo             bool     `json:"echo"`
	Logprobs         int      `json:"logprobs"`
	Stop             []string `json:"stop"`
	Stream           bool     `json:"stream"`
}

func init() {
	validation.AddCustomFunc("ValidateTemperature", validateTemperature)
	validation.AddCustomFunc("ValidatePenalty", validatePenalty)
}

var validateTemperature validation.CustomFunc = func(v *validation.Validation, obj interface{}, key string) {
	temperature, ok := obj.(float32)

	if !ok {
		return
	}

	if temperature < 0.0 || temperature > 1.0 {
		v.AddError(key, "must in 0~1")
	}
}

var validatePenalty validation.CustomFunc = func(v *validation.Validation, obj interface{}, key string) {
	temperature, ok := obj.(float32)

	if !ok {
		return
	}

	if temperature < 0.0 || temperature > 2.0 {
		v.AddError(key, "must in 0~2")
	}
}
