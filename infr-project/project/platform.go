package project

type ProjectPlatform struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// ref platform.PlatformProject
	ProjectId string `json:"project_id"`
}
