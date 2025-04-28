package evaluation

type EvaluationResult struct {
	Metrics map[string]IEvaluationMetric
}

func NewEvaluationResult(metrics []IEvaluationMetric) *EvaluationResult {
	e := &EvaluationResult{
		Metrics: make(map[string]IEvaluationMetric),
	}

	for _, v := range metrics {
		e.Metrics[v.GetName()] = v
	}

	return e
}

func (r *EvaluationResult) AddOrUpdateContextInAllMetrics(values []EvaluationContext) {
	for _, v := range r.Metrics {
		v.AddOrUpdateContext(values)
	}
}

func (r *EvaluationResult) AddDiagnosticsToAllMetrics(diagnostics []EvaluationDiagnostic) {
	for _, v := range r.Metrics {
		v.AddDiagnostics(diagnostics)
	}
}

func (r *EvaluationResult) ContainsDiagnostics(predicate func(EvaluationDiagnostic) bool) bool {
	if predicate == nil {
		return false
	}
	for _, v := range r.Metrics {
		if v.ContainsDiagnostics(predicate) {
			return true
		}
	}
	return false
}

func (r *EvaluationResult) Interpret(interpretationProvider func(IEvaluationMetric) *EvaluationMetricInterpretation) {
	for _, v := range r.Metrics {
		if i := interpretationProvider(v); i != nil {
			v.SetInterpretation(i)
		}
	}
}

func (r *EvaluationResult) AddOrUpdateMetadataInAllMetrics(name string, value string) {
	for _, metric := range r.Metrics {
		metric.AddOrUpdateMetadata(name, value)
	}
}
