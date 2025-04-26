package evaluation

type BooleanMetric EvaluationMetricT[bool]

func NewBooleanMetric(name string, value bool, reason *string) BooleanMetric {
	return BooleanMetric{
		EvaluationMetric: NewEvaluationMetric(name, reason),
		Value:			value,
	}
}