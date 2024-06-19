package project

// ProjectState is the interface for webhook states.
type ProjectState interface {
	privateProjectState() // Prevents external implementation
	String() string
}

// projectState is the underlying implementation for ProjectState.
type projectState string

// privateWebhookState makes webhookState implement ProjectState.
func (c projectState) privateProjectState() {}

// String makes webhookState implement WebhookState.
func (c projectState) String() string {
	return string(c)
}

// Constants for the different webhook states.
const (
	ProjectPreparing  projectState = "preparing"
	ProjectProcessing projectState = "processing"
	ProjectFinished   projectState = "finished"
)

func GetProjectState(rType string) ProjectState {
	switch rType {
	case "preparing":
		return ProjectPreparing
	case "processing":
		return ProjectProcessing
	case "finished":
		return ProjectFinished
	default:
		return ProjectPreparing
	}
}
