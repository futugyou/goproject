package vercel

import (
	"strings"
)

type AccessGroupService service

func (v *AccessGroupService) CreateAccessGroup(slug string, teamId string, info AccessGroupRequest) (*AccessGroupInfo, error) {
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
	result := &AccessGroupInfo{}
	err := v.client.http.Post(path, info, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AccessGroupService) DeleteAccessGroup(idOrName string, slug string, teamId string) (*string, error) {
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
	err := v.client.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *AccessGroupService) ListMembersOfAccessGroup(idOrName string, slug string, teamId string) (*ListMembersResponse, error) {
	path := "/v1/access-groups/" + idOrName + "/members"
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
	result := &ListMembersResponse{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AccessGroupService) ListProjectsOfAccessGroup(idOrName string, slug string, teamId string) (*ListProjectsResponse, error) {
	path := "/v1/access-groups/" + idOrName + "/projects"
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
	result := &ListProjectsResponse{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AccessGroupService) ListAccessGroup(projectId string, search string, slug string, teamId string) ([]AccessGroupInfo, error) {
	path := "/v1/access-groups/"
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
	if len(projectId) > 0 {
		if strings.Contains(path, "?") {
			path += ("&projectId=" + projectId)
		} else {
			path += ("?projectId=" + projectId)
		}
	}
	if len(search) > 0 {
		if strings.Contains(path, "?") {
			path += ("&search=" + search)
		} else {
			path += ("?search=" + search)
		}
	}
	result := []AccessGroupInfo{}
	err := v.client.http.Get(path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AccessGroupService) GetAccessGroup(idOrName string, slug string, teamId string) (*AccessGroupInfo, error) {
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
	result := &AccessGroupInfo{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AccessGroupService) UpdateAccessGroup(idOrName string, slug string, teamId string, info AccessGroupRequest) (*AccessGroupInfo, error) {
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
	result := &AccessGroupInfo{}
	err := v.client.http.Post(path, info, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type AccessGroupRequest struct {
	Name            string               `json:"name"`
	MembersToAdd    []string             `json:"membersToAdd"`
	MembersToRemove []string             `json:"membersToRemove,omitempty"`
	Projects        []AccessGroupProject `json:"projects"`
}

type AccessGroupProject struct {
	ProjectID          string `json:"projectId"`
	Role               string `json:"role"`
	Framework          string `json:"framework"`
	LatestDeploymentId string `json:"latestDeploymentId"`
	Name               string `json:"name"`
}

type AccessGroupInfo struct {
	AccessGroupId string       `json:"accessGroupId"`
	CreatedAt     string       `json:"createdAt"`
	MembersCount  int          `json:"membersCount"`
	Name          string       `json:"Name"`
	ProjectsCount int          `json:"projectsCount"`
	TeamId        string       `json:"teamId"`
	UpdatedAt     string       `json:"updatedAt"`
	Error         *VercelError `json:"error"`
}

type ListMembersResponse struct {
	Members    []MemberInfo `json:"members"`
	Pagination Pagination   `json:"pagination"`
	Error      *VercelError `json:"error"`
}

type ListProjectsResponse struct {
	Projects   []AccessGroupProjectInfo `json:"projects"`
	Pagination Pagination               `json:"pagination"`
	Error      *VercelError             `json:"error"`
}

type AccessGroupProjectInfo struct {
	CreatedAt string             `json:"createdAt"`
	Project   AccessGroupProject `json:"project"`
	ProjectId string             `json:"projectId"`
	Role      string             `json:"role"`
	UpdatedAt string             `json:"updatedAt"`
}
