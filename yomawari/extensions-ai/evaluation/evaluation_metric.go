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

func NewEvaluationMetric(name string, reason *string) EvaluationMetric {
	return EvaluationMetric{
		Name:   name,
		Reason: reason,
	}
}

func (e *EvaluationMetric) AddOrUpdateContext(name string, value []contents.IAIContent) {
	e.Context[name] = value
}

func (e *EvaluationMetric) AddDiagnostics(diagnostics []EvaluationDiagnostic) {
	if len(diagnostics) > 0 {
		e.Diagnostics = append(e.Diagnostics, diagnostics...)
	}
}

func (e *EvaluationMetric) AddOrUpdateMetadata(name string, value string) {
	e.Metadata[name] = value
}

func (e *EvaluationMetric) ContainsDiagnostics(predicate func(EvaluationDiagnostic) bool) bool {
	if e == nil || e.Diagnostics == nil || predicate == nil {
		return false
	}

	for _, v := range e.Diagnostics {
		if predicate(v) {
			return true
		}
	}
	return false
}
