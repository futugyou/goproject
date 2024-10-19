package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type MemberService service

func (v *MemberService) AddMember(ctx context.Context, idOrName string, slug string, teamId string, req AddMemberRequest) (*OperateMemberResponse, error) {
	path := fmt.Sprintf("/v1/projects/%s/members", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &OperateMemberResponse{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *MemberService) ListProjectMember(ctx context.Context, idOrName string, slug string, teamId string, search string, limit string, since string, until string) (*ListProjectMemberrResponse, error) {
	path := fmt.Sprintf("/v1/projects/%s/members", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(limit) > 0 {
		queryParams.Add("limit", limit)
	}
	if len(since) > 0 {
		queryParams.Add("since", since)
	}
	if len(until) > 0 {
		queryParams.Add("until", until)
	}
	if len(search) > 0 {
		queryParams.Add("search", search)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &ListProjectMemberrResponse{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *MemberService) RemoveProjectMember(ctx context.Context, idOrName string, uid string, slug string, teamId string) (*OperateMemberResponse, error) {
	path := fmt.Sprintf("/v1/projects/%s/members/%s", idOrName, uid)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &OperateMemberResponse{}
	err := v.client.http.Delete(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type AddMemberRequest struct {
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
	Uid      string `json:"uid,omitempty"`
	Username string `json:"username,omitempty"`
}

type OperateMemberResponse struct {
	Id    string       `json:"id,omitempty"`
	Error *VercelError `json:"error,omitempty"`
}

type ListProjectMemberrResponse struct {
	Members    []MemberInfo `json:"members,omitempty"`
	Pagination Pagination   `json:"pagination,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}

type MemberInfo struct {
	Avatar              string `json:"avatar"`
	CreatedAt           string `json:"createdAt"`
	ComputedProjectRole string `json:"computedProjectRole"`
	Email               int    `json:"email"`
	Name                string `json:"name"`
	TeamRole            string `json:"teamRole"`
	Role                string `json:"role"`
	Uid                 string `json:"uid"`
	Username            string `json:"username"`
}
