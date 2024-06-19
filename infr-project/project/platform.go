package project

// json and bson tag is necessary, we have not custom serialization
type ProjectPlatform struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	// ref platform.PlatformProject
	ProjectId string `json:"project_id" bson:"project_id"`
}
