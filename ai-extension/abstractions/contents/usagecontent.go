package contents

import (
	"github.com/futugyou/ai-extension/abstractions"
)

// UsageContent represents content that holds usage information.
type UsageContent struct {
	AIContent `json:",inline"`
	Details   abstractions.UsageDetails `json:"details,omitempty"`
}
