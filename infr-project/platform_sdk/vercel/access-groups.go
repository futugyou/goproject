package vercel

import (
	"strings"
)

func (v *VercelClient) CreateAccessGroup(slug string, teamId string, info AccessGroupInfo) (*CreateAccessGroupResponse, error) {
	path := "/v1/access-groups"
	if len(slug) > 0 {
		path += ("?slug=" + slug)
	}
	if len(teamId) > 0 {
		if strings.Contains(path, "?") {
			path += ("&teamId=" + teamId)
		} else {
			path += ("?teamId=" + teamId)
		}
	}
	result := &CreateAccessGroupResponse{}
	err := v.http.Post(path, info, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) DeleteAccessGroup(idOrName string, slug string, teamId string) (*string, error) {
	path := "/v1/access-groups/" + idOrName
	if len(slug) > 0 {
		path += ("?slug=" + slug)
	}
	if len(teamId) > 0 {
		if strings.Contains(path, "?") {
			path += ("&teamId=" + teamId)
		} else {
			path += ("?teamId=" + teamId)
		}
	}
	result := ""
	err := v.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

type AccessGroupInfo struct {
	Name         string    `json:"name"`
	MembersToAdd string    `json:"membersToAdd"`
	Projects     []Project `json:"projects"`
}

type Project struct {
	ProjectID string `json:"projectId"`
	Role      string `json:"role"`
}

type CreateAccessGroupResponse struct {
	AccessGroupId string       `json:"accessGroupId"`
	CreatedAt     string       `json:"createdAt"`
	MembersCount  int          `json:"membersCount"`
	Name          string       `json:"Name"`
	ProjectsCount int          `json:"projectsCount"`
	TeamId        string       `json:"teamId"`
	UpdatedAt     string       `json:"updatedAt"`
	Error         *VercelError `json:"error"`
}
