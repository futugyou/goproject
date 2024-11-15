package project

// json and bson tag is necessary, we have not custom serialization
type ProjectDesign struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	// ref resource.Resource
	Resources []ProjectResource `json:"resources" bson:"resources"`
}

type ProjectResource struct {
	Name string `json:"name" bson:"name"`
	// ref resource.Resource
	ResourceId string `json:"resource_id" bson:"resource_id"`
}
