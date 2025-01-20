package project

type ProjectPlatform struct {
	Name        string `bson:"name"`
	Description string `bson:"description"`
	// ref platform.PlatformProject
	PlatformId string `bson:"platform_id"`
	ProjectId  string `bson:"project_id"`
}
