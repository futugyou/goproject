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
