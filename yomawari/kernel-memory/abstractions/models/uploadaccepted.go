package models

type UploadAccepted struct {
	Index      string `json:"index"`
	DocumentId string `json:"documentId"`
	Message    string `json:"message"`
}
