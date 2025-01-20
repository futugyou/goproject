package project

type ProjectDesign struct {
	Name        string `bson:"name"`
	Description string `bson:"description"`
	// ref resource.Resource
	Resources []ProjectResource `bson:"resources"`
}

type ProjectResource struct {
	Name string `bson:"name"`
	// ref resource.Resource
	ResourceId string `bson:"resource_id"`
}
