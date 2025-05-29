package protocol

import "strings"

type ModelPreferences struct {
	CostPriority         *float32    `json:"costPriority"`
	Hints                []ModelHint `json:"hints"`
	SpeedPriority        *float32    `json:"speedPriority"`
	IntelligencePriority *float32    `json:"intelligencePriority"`
}

func (m *ModelPreferences) Validate() (string, bool) {
	valid := true
	errorMessage := ""
	errors := []string{}

	if m.CostPriority == nil || *m.CostPriority < 0 || *m.CostPriority > 1 {
		errors = append(errors, "CostPriority must be between 0 and 1")
		valid = false
	}

	if m.SpeedPriority == nil || *m.SpeedPriority < 0 || *m.SpeedPriority > 1 {
		errors = append(errors, "SpeedPriority must be between 0 and 1")
		valid = false
	}

	if m.IntelligencePriority == nil || *m.IntelligencePriority < 0 || *m.IntelligencePriority > 1 {
		errors = append(errors, "IntelligencePriority must be between 0 and 1")
		valid = false
	}

	if !valid {
		errorMessage = strings.Join(errors, ",")
	}

	return errorMessage, valid
}
