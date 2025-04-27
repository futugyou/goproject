package evaluation

import (
	"github.com/futugyou/yomawari/extensions-ai/abstractions/contents"
)

type IEvaluationMetric interface {
	AddOrUpdateContext(name string, value []contents.IAIContent)
	AddDiagnostics(diagnostics []EvaluationDiagnostic)
	AddOrUpdateMetadata(name string, value string)
	ContainsDiagnostics(predicate func(EvaluationDiagnostic) bool) bool
	GetName() string
	SetName(string)
	GetInterpretation() *EvaluationMetricInterpretation
	SetInterpretation(*EvaluationMetricInterpretation)
}

var _ IEvaluationMetric = (*EvaluationMetric)(nil)

type EvaluationMetric struct {
	name           string
	Reason         *string
	interpretation *EvaluationMetricInterpretation
	Context        map[string][]contents.IAIContent
	Metadata       map[string]string
	Diagnostics    []EvaluationDiagnostic
}

func NewEvaluationMetric(name string, reason *string) EvaluationMetric {
	return EvaluationMetric{
		name:   name,
		Reason: reason,
	}
}

func (e *EvaluationMetric) AddOrUpdateContext(name string, value []contents.IAIContent) {
	e.Context[name] = value
}

func (e *EvaluationMetric) GetName() string {
	return e.name
}

func (e *EvaluationMetric) SetName(n string) {
	e.name = n
}

func (e *EvaluationMetric) GetInterpretation() *EvaluationMetricInterpretation {
	return e.interpretation
}

func (e *EvaluationMetric) SetInterpretation(a *EvaluationMetricInterpretation) {
	e.interpretation = a
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
