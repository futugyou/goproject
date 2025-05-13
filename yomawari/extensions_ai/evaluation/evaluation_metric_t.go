package evaluation

var _ IEvaluationMetric = (*EvaluationMetricT[any])(nil)

type EvaluationMetricT[T any] struct {
	EvaluationMetric
	Value T
}

func NewEvaluationMetricT[T any](name string, value T, reason *string) EvaluationMetricT[T] {
	return EvaluationMetricT[T]{
		EvaluationMetric: NewEvaluationMetric(name, reason),
		Value:            value,
	}
}
