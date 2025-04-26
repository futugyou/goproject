package evaluation

type NumericMetric EvaluationMetricT[float64]

func NewNumericMetric(name string, value float64, reason *string) NumericMetric {
	return NumericMetric{
		EvaluationMetric: NewEvaluationMetric(name, reason),
		Value:            value,
	}
}
