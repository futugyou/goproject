package domain

type ProjectDesign struct {
	Name            string `bson:"name"`
	Description     string `bson:"description"`
	ResourceID      string `bson:"resource_id"`
	ResourceVersion int    `bson:"resource_version"`
}
