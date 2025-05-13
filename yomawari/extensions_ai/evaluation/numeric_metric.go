package evaluation

import "strconv"

var _ IEvaluationMetric = (*NumericMetric)(nil)

type NumericMetric EvaluationMetricT[float64]

func NewNumericMetric(name string, value float64, reason *string) NumericMetric {
	return NumericMetric{
		EvaluationMetric: NewEvaluationMetric(name, reason),
		Value:            value,
	}
}

func (metric *NumericMetric) InterpretScore() EvaluationMetricInterpretation {
	const MinimumPassingScore = 4.0

	var rating EvaluationRating

	if metric == nil {
		rating = EvaluationRatingInconclusive
	} else {
		v := metric.Value
		switch {
		case v > 5.0:
			rating = EvaluationRatingInconclusive
		case v > 4.0 && v <= 5.0:
			rating = EvaluationRatingExceptional
		case v > 3.0 && v <= 4.0:
			rating = EvaluationRatingGood
		case v > 2.0 && v <= 3.0:
			rating = EvaluationRatingAverage
		case v > 1.0 && v <= 2.0:
			rating = EvaluationRatingPoor
		case v > 0.0 && v <= 1.0:
			rating = EvaluationRatingUnacceptable
		case v <= 0.0:
			rating = EvaluationRatingInconclusive
		default:
			rating = EvaluationRatingInconclusive
		}
	}

	if metric != nil && metric.Value < MinimumPassingScore {
		reason := metric.GetName() + " is less than " + formatFloat(MinimumPassingScore) + "."
		return EvaluationMetricInterpretation{
			Rating: rating,
			Failed: true,
			Reason: &reason,
		}
	}

	return EvaluationMetricInterpretation{
		Rating: rating,
	}
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
