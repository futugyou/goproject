package abstractions

type UsageDetails struct {
	InputTokenCount      *int64           `json:"inputTokenCount,omitempty"`
	OutputTokenCount     *int64           `json:"outputTokenCount,omitempty"`
	TotalTokenCount      *int64           `json:"totalTokenCount,omitempty"`
	AdditionalProperties map[string]int64 `json:"additionalProperties,omitempty"`
}

func (u *UsageDetails) AddUsageDetails(usage UsageDetails) {
	if u == nil {
		u = &UsageDetails{
			InputTokenCount:      new(int64),
			OutputTokenCount:     new(int64),
			TotalTokenCount:      new(int64),
			AdditionalProperties: map[string]int64{},
		}
	}

	if usage.InputTokenCount != nil {
		if u.InputTokenCount == nil {
			u.InputTokenCount = new(int64)
		}
		*u.InputTokenCount += *usage.InputTokenCount
	}

	if usage.OutputTokenCount != nil {
		if u.OutputTokenCount == nil {
			u.OutputTokenCount = new(int64)
		}
		*u.OutputTokenCount += *usage.OutputTokenCount
	}

	if usage.TotalTokenCount != nil {
		if u.TotalTokenCount == nil {
			u.TotalTokenCount = new(int64)
		}
		*u.TotalTokenCount += *usage.TotalTokenCount
	}

	if u.AdditionalProperties == nil {
		u.AdditionalProperties = make(map[string]int64)
	}

	for key, value := range usage.AdditionalProperties {
		if existingValue, exists := u.AdditionalProperties[key]; exists {
			u.AdditionalProperties[key] = existingValue + value
		} else {
			u.AdditionalProperties[key] = value
		}
	}
}
