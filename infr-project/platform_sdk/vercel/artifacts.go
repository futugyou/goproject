package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CheckArtifactExists(hash string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v8/artifacts/%s", hash)
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
	result := ""
	err := v.http.Get(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}
