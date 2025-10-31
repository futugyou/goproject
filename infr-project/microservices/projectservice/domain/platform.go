package domain

type ProjectPlatform struct {
	Name        string `bson:"name"`
	Description string `bson:"description"`
	// ref platform.PlatformProject
	PlatformID string `bson:"platform_id"`
	ProjectID  string `bson:"project_id"`
}
