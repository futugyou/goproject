package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) AddMember(idOrName string, slug string, teamId string, req AddMemberRequest) (*AddMemberResponse, error) {
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
	result := &AddMemberResponse{}
	err := v.http.Post(path, req, result)

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

type AddMemberResponse struct {
	Id    string       `json:"id,omitempty"`
	Error *VercelError `json:"error,omitempty"`
}
