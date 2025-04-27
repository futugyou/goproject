package evaluation

var _ IEvaluationMetric = (*StringMetric)(nil)

type StringMetric EvaluationMetricT[string]

func NewStringMetric(name string, value string, reason *string) StringMetric {
	return StringMetric{
		EvaluationMetric: NewEvaluationMetric(name, reason),
		Value:            value,
	}
}
