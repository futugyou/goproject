package evaluation

type EvaluationMetricInterpretation struct {
	Rating EvaluationRating
	Failed bool
	Reason *string
}

func NewEvaluationMetricInterpretation(rating *EvaluationRating, failed *bool, reason *string) *EvaluationMetricInterpretation {
	e := &EvaluationMetricInterpretation{Rating: EvaluationRatingUnknown, Failed: false, Reason: reason}
	if rating != nil {
		e.Rating = *rating
	}
	if failed != nil {
		e.Failed = *failed
	}
	return e
}
