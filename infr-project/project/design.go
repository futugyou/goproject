package project

type ProjectDesign struct {
	Name            string `bson:"name"`
	Description     string `bson:"description"`
	ResourceId      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
}
