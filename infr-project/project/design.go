package project

// json and bson tag is necessary, we have not custom serialization
type ProjectDesign struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	// ref resource.Resource
	Resources []string `json:"resources" bson:"resources"`
}
