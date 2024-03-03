package openapi

type OpenAPIConfig struct {
	SpceVersion string            `json:"spce_version"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	APIVersion  string            `json:"apiVersion"`
	ModelFolder string            `json:"model"`
	OutputPath  string            `json:"output"`
	OutputType  string            `json:"type"`
	APIConfigs  []OperationConfig `json:"apis"`
}

type OperationConfig struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Request     string `json:"request"`
	Response    string `json:"response"`
	Description string `json:"description"`
}

func (m *OpenAPIConfig) Check() error {
	if len(m.Title) == 0 {
		m.Title = "this is openapi tile"
	}
	if len(m.Description) == 0 {
		m.Description = "this is openapi description"
	}
	if len(m.APIVersion) == 0 {
		m.APIVersion = "0.0.0"
	}
	if len(m.ModelFolder) == 0 {
		m.ModelFolder = "./"
	}
	if len(m.OutputPath) == 0 {
		m.OutputPath = "./openapi.json"
	}
	if len(m.OutputType) == 0 {
		m.OutputType = "json"
	}
	if len(m.SpceVersion) == 0 {
		m.SpceVersion = "3.1.3"
	}
	return nil
}
