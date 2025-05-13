package evaluation

type EvaluationDiagnostic struct {
	Severity EvaluationDiagnosticSeverity
	Message  string
}

func EvaluationDiagnosticInformational(message string) EvaluationDiagnostic {
	return EvaluationDiagnostic{
		Severity: EvaluationDiagnosticSeverityInformational,
		Message:  message,
	}
}

func EvaluationDiagnosticWarning(message string) EvaluationDiagnostic {
	return EvaluationDiagnostic{
		Severity: EvaluationDiagnosticSeverityWarning,
		Message:  message,
	}
}

func EvaluationDiagnosticError(message string) EvaluationDiagnostic {
	return EvaluationDiagnostic{
		Severity: EvaluationDiagnosticSeverityError,
		Message:  message,
	}
}
