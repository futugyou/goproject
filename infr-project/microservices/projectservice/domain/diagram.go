package domain

type Diagram struct {
	Type            string `json:"type"`
	Description     string `json:"description,omitempty"`
	ResourceID      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
}
