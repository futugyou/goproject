package models

import "time"

type TokenUsage struct {
	Timestamp              time.Time `json:"timestamp"`
	ServiceType            *string   `json:"serviceType,omitempty"`
	ModelType              *string   `json:"modelType,omitempty"`
	ModelName              *string   `json:"modelName,omitempty"`
	TokenizerTokensIn      *int      `json:"tokenizerTokensIn,omitempty"`
	TokenizerTokensOut     *int      `json:"tokenizerTokensOut,omitempty"`
	ServiceTokensIn        *int      `json:"serviceTokensIn,omitempty"`
	ServiceTokensOut       *int      `json:"serviceTokensOut,omitempty"`
	ServiceReasoningTokens *int      `json:"serviceReasoningTokens,omitempty"`
}

// Merge combines the fields from another TokenUsage into the current instance
func (t *TokenUsage) Merge(input *TokenUsage) {
	if input == nil || t == nil {
		return
	}

	t.Timestamp = input.Timestamp
	t.ServiceType = input.ServiceType
	t.ModelType = input.ModelType
	t.ModelName = input.ModelName

	// Merge the token counts by adding them together, considering nil values
	if input.TokenizerTokensIn != nil {
		if t.TokenizerTokensIn == nil {
			t.TokenizerTokensIn = new(int)
		}
		*t.TokenizerTokensIn += *input.TokenizerTokensIn
	}

	if input.TokenizerTokensOut != nil {
		if t.TokenizerTokensOut == nil {
			t.TokenizerTokensOut = new(int)
		}
		*t.TokenizerTokensOut += *input.TokenizerTokensOut
	}

	if input.ServiceTokensIn != nil {
		if t.ServiceTokensIn == nil {
			t.ServiceTokensIn = new(int)
		}
		*t.ServiceTokensIn += *input.ServiceTokensIn
	}

	if input.ServiceTokensOut != nil {
		if t.ServiceTokensOut == nil {
			t.ServiceTokensOut = new(int)
		}
		*t.ServiceTokensOut += *input.ServiceTokensOut
	}

	if input.ServiceReasoningTokens != nil {
		if t.ServiceReasoningTokens == nil {
			t.ServiceReasoningTokens = new(int)
		}
		*t.ServiceReasoningTokens += *input.ServiceReasoningTokens
	}
}
