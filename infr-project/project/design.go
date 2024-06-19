package project

type ProjectDesign struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// ref resource.Resource
	Resources []string `json:"resources"`
}
