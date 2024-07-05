package vercel

import (
	"fmt"
	"strings"
)

func (v *VercelClient) AssignAlias(id string, slug string, teamId string, info AssignAliasRequest) (*AssignAliasResponse, error) {
	path := fmt.Sprintf("/v2/deployments/%s/aliases", id)
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
	result := &AssignAliasResponse{}
	err := v.http.Post(path, info, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type AssignAliasRequest struct {
	Alias    string       `json:"alias"`
	Redirect string       `json:"redirect,omitempty"`
	Error    *VercelError `json:"error"`
}

type AssignAliasResponse struct {
	Alias           string       `json:"alias"`
	Created         string       `json:"created"`
	Uid             string       `json:"uid"`
	OldDeploymentId string       `json:"oldDeploymentId"`
	Error           *VercelError `json:"error"`
}
