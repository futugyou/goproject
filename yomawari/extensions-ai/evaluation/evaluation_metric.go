package evaluation

import (
	"github.com/futugyou/yomawari/extensions-ai/abstractions/contents"
)

type EvaluationMetric struct {
	Name           string
	Reason         *string
	Interpretation *EvaluationMetricInterpretation
	Context        map[string][]contents.IAIContent
	Metadata       map[string]string
	Diagnostics    []EvaluationDiagnostic
}
