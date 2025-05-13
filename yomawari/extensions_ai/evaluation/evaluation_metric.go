package evaluation

import (
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

type IEvaluationMetric interface {
	AddOrUpdateContext(values []EvaluationContext)
	AddDiagnostics(diagnostics []EvaluationDiagnostic)
	AddOrUpdateMetadata(name string, value string)
	ContainsDiagnostics(predicate func(EvaluationDiagnostic) bool) bool
	GetName() string
	SetName(name string)
	GetInterpretation() *EvaluationMetricInterpretation
	SetInterpretation(*EvaluationMetricInterpretation)
	AddOrUpdateChatMetadata(response chatcompletion.ChatResponse, duration float64)
}

var _ IEvaluationMetric = (*EvaluationMetric)(nil)

type EvaluationMetric struct {
	name           string
	Reason         *string
	interpretation *EvaluationMetricInterpretation
	Context        map[string]EvaluationContext
	Metadata       map[string]string
	Diagnostics    []EvaluationDiagnostic
}

// AddOrUpdateChatMetadata implements IEvaluationMetric.
func (metric *EvaluationMetric) AddOrUpdateChatMetadata(evaluationResponse chatcompletion.ChatResponse, duration float64) {
	if evaluationResponse.ModelId != nil {
		metric.AddOrUpdateMetadata("evaluation-model-used", *evaluationResponse.ModelId)
	}

	if evaluationResponse.Usage != nil {
		if evaluationResponse.Usage.InputTokenCount != nil {
			metric.AddOrUpdateMetadata("evaluation-input-tokens-used", fmt.Sprintf("%d", *evaluationResponse.Usage.InputTokenCount))
		}

		if evaluationResponse.Usage.OutputTokenCount != nil {
			metric.AddOrUpdateMetadata("evaluation-output-tokens-used", fmt.Sprintf("%d", *evaluationResponse.Usage.OutputTokenCount))
		}

		if evaluationResponse.Usage.TotalTokenCount != nil {
			metric.AddOrUpdateMetadata("evaluation-total-tokens-used", fmt.Sprintf("%d", *evaluationResponse.Usage.TotalTokenCount))
		}
	}
	if duration > 0 {
		duration := fmt.Sprintf("%f", duration)
		metric.AddOrUpdateMetadata("evaluation-duration", duration)
	}
}

func NewEvaluationMetric(name string, reason *string) EvaluationMetric {
	return EvaluationMetric{
		name:   name,
		Reason: reason,
	}
}

func (e *EvaluationMetric) AddOrUpdateContext(values []EvaluationContext) {
	for _, v := range values {
		e.Context[v.GetName()] = v
	}
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
