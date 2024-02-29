package openapi

type OpenAPIConfig struct {
	Title       string
	Description string
	Version     string
	ModelFolder string
	OutputPath  string
	OutputType  string
}

func (m *OpenAPIConfig) Check() error {
	if len(m.Title) == 0 {
		m.Title = ""
	}
	if len(m.Description) == 0 {
		m.Description = ""
	}
	if len(m.Version) == 0 {
		m.Version = ""
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
	return nil
}
