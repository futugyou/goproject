package domain

// WebhookState is the interface for webhook states.
type WebhookState interface {
	privateWebhookState() // Prevents external implementation
	String() string
}

// webhookState is the underlying implementation for WebhookState.
type webhookState string

// privateWebhookState makes webhookState implement WebhookState.
func (c webhookState) privateWebhookState() {}

// String makes webhookState implement WebhookState.
func (c webhookState) String() string {
	return string(c)
}

// Constants for the different webhook states.
const (
	WebhookCreating webhookState = "Creating"
	WebhookReady    webhookState = "Ready"
	WebhookFailed   webhookState = "Failed"
	WebhookDeleted  webhookState = "Deleted"
)

func GetWebhookState(rType string) WebhookState {
	switch rType {
	case "Creating":
		return WebhookCreating
	case "Ready":
		return WebhookReady
	case "Failed":
		return WebhookFailed
	case "Deleted":
		return WebhookDeleted
	default:
		return WebhookCreating
	}
}
