package vercel

import "log"

func (v *VercelClient) GetProjects() string {
	path := "/v9/projects"
	result := "[]"
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (v *VercelClient) GetProjectEnv(project string) string {
	path := "/v9/projects/" + project + "/env"
	result := "[]"
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (v *VercelClient) GetProject(project string, slug string, teamId string) string {
	path := "/v9/projects/" + project + "?slug=" + slug + "&teamId=" + teamId
	result := ""
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (v *VercelClient) CreateProject(name string, slug string, teamId string) string {
	path := "/v10/projects?slug=" + slug + "&teamId=" + teamId
	// TODO: struct
	result := ""
	info := ProjectInfo{
		Name: name,
	}
	err := v.http.Post(path, info, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

type ProjectInfo struct {
	Name                                 string                `json:"name"`
	BuildCommand                         string                `json:"buildCommand,omitempty"`
	CommandForIgnoringBuildStep          string                `json:"commandForIgnoringBuildStep,omitempty"`
	DevCommand                           string                `json:"devCommand,omitempty"`
	EnvironmentVariables                 []EnvironmentVariable `json:"environmentVariables,omitempty"`
	Framework                            string                `json:"framework,omitempty"`
	GitRepository                        GitRepository         `json:"gitRepository,omitempty"`
	InstallCommand                       string                `json:"installCommand,omitempty"`
	OutputDirectory                      string                `json:"outputDirectory,omitempty"`
	PublicSource                         bool                  `json:"publicSource,omitempty"`
	RootDirectory                        string                `json:"rootDirectory,omitempty"`
	ServerlessFunctionRegion             string                `json:"serverlessFunctionRegion,omitempty"`
	ServerlessFunctionZeroConfigFailover string                `json:"serverlessFunctionZeroConfigFailover,omitempty"`
	SkipGitConnectDuringLink             bool                  `json:"skipGitConnectDuringLink,omitempty"`
}

type EnvironmentVariable struct {
	Key       string `json:"key,omitempty"`
	Target    string `json:"target,omitempty"`
	GitBranch string `json:"gitBranch,omitempty"`
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
}

type GitRepository struct {
	Repo string `json:"repo,omitempty"`
	Type string `json:"type,omitempty"`
}
