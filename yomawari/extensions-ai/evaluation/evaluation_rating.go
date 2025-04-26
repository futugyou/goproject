package evaluation

type EvaluationRating string

const (
	EvaluationRatingUnknown      EvaluationRating = "Unknown"
	EvaluationRatingInconclusive EvaluationRating = "Inconclusive"
	EvaluationRatingExceptional  EvaluationRating = "Exceptional"
	EvaluationRatingGood         EvaluationRating = "Good"
	EvaluationRatingAverage      EvaluationRating = "Average"
	EvaluationRatingPoor         EvaluationRating = "Poor"
	EvaluationRatingUnacceptable EvaluationRating = "Unacceptable"
)
