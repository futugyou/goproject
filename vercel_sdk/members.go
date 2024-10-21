package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type MemberService service

type AddMemberRequest struct {
	Role             string `json:"role,omitempty"`
	Email            string `json:"email,omitempty"`
	Uid              string `json:"uid,omitempty"`
	Username         string `json:"username,omitempty"`
	IdOrName         string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *MemberService) AddMember(ctx context.Context, request AddMemberRequest) (*OperateMemberResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/members", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &OperateMemberResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListProjectMemberParameter struct {
	IdOrName         string `json:"-"`
	Limit            *string
	Search           *string
	Since            *string
	Until            *string
	BaseUrlParameter `json:"-"`
}

func (v *MemberService) ListProjectMember(ctx context.Context, request ListProjectMemberParameter) (*ListProjectMemberrResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/members", request.IdOrName),
	}
	params := request.GetUrlValues()
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	if request.Search != nil {
		params.Add("search", *request.Search)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListProjectMemberrResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type RemoveProjectMemberRequest struct {
	IdOrName         string
	Uid              string
	BaseUrlParameter `json:"-"`
}

func (v *MemberService) RemoveProjectMember(ctx context.Context, request RemoveProjectMemberRequest) (*OperateMemberResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/members/%s", request.IdOrName, request.Uid),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &OperateMemberResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
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
