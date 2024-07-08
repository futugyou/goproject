package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) DeleteIntegrationConfiguration(id string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/integrations/configuration/%s", id)
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
	err := v.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}
