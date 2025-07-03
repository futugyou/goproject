package project

type Diagram struct {
	Type            string `json:"type"`
	Description     string `json:"description,omitempty"`
	ResourceId      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
}
