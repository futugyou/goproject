package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type AccessGroupService service

func (v *AccessGroupService) CreateAccessGroup(ctx context.Context, request AccessGroupRequest) (*AccessGroupInfo, error) {
	u := &url.URL{
		Path: "/v1/access-groups",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &AccessGroupInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}

	return response, nil
}

type DeleteAccessGroupRequest struct {
	IdOrName         string `json:"idOrName"`
	BaseUrlParameter `json:"-"`
}

func (v *AccessGroupService) DeleteAccessGroup(ctx context.Context, request DeleteAccessGroupRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/access-groups/%s", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}

	return &response, nil
}

type ListMembersOfAccessGroupParameter struct {
	IdOrName         string `json:"idOrName"`
	BaseUrlParameter `json:"-"`
}

func (v *AccessGroupService) ListMembersOfAccessGroup(ctx context.Context, request ListMembersOfAccessGroupParameter) (*ListMembersResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/access-groups/%s/members", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListMembersResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListProjectsOfAccessGroupParameter struct {
	IdOrName         string `json:"idOrName"`
	BaseUrlParameter `json:"-"`
}

func (v *AccessGroupService) ListProjectsOfAccessGroup(ctx context.Context, request ListProjectsOfAccessGroupParameter) (*ListProjectsResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/access-groups/%s/projects", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListProjectsResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListAccessGroupParameter struct {
	ProjectId        *string `json:"projectId"`
	Search           *string `json:"search"`
	BaseUrlParameter `json:"-"`
}

func (v *AccessGroupService) ListAccessGroup(ctx context.Context, request ListAccessGroupParameter) ([]AccessGroupInfo, error) {
	u := &url.URL{
		Path: "/v1/access-groups/",
	}
	params := request.GetUrlValues()
	if request.Search != nil {
		params.Add("search", *request.Search)
	}
	if request.ProjectId != nil {
		params.Add("projectId", *request.ProjectId)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []AccessGroupInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return response, nil
}

type GetAccessGroupParameter struct {
	IdOrName         string `json:"idOrName"`
	BaseUrlParameter `json:"-"`
}

func (v *AccessGroupService) GetAccessGroup(ctx context.Context, request GetAccessGroupParameter) (*AccessGroupInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/access-groups/%s", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &AccessGroupInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (v *AccessGroupService) UpdateAccessGroup(ctx context.Context, request AccessGroupRequest) (*AccessGroupInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/access-groups/%s", *request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &AccessGroupInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}

	return response, nil
}

type AccessGroupRequest struct {
	Name             string               `json:"name"`
	MembersToAdd     []string             `json:"membersToAdd"`
	MembersToRemove  []string             `json:"membersToRemove,omitempty"`
	Projects         []AccessGroupProject `json:"projects"`
	IdOrName         *string              `json:"-"`
	BaseUrlParameter `json:"-"`
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
