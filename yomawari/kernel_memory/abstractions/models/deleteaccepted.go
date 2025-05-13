package models

type DeleteAccepted struct {
	Index      string  `json:"index"`
	DocumentId *string `json:"documentId,omitempty"`
	Message    string  `json:"message"`
}
