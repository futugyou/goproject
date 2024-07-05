package vercel

import (
	"net/url"
)

func (v *VercelClient) CreateAuthToken(slug string, teamId string, name string) (*CreateAuthTokenResponse, error) {
	path := "/v3/user/tokens"
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
	result := &CreateAuthTokenResponse{}
	request := struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		Name: name,
		Type: "oauth2-token",
	}
	err := v.http.Post(path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateAuthTokenResponse struct {
	BearerToken string       `json:"bearerToken"`
	Token       TokenInfo    `json:"token"`
	Error       *VercelError `json:"error"`
}

type TokenInfo struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Origin    string  `json:"origin"`
	Scopes    []Scope `json:"scopes"`
	ExpiresAt int     `json:"expiresAt"`
	ActiveAt  int     `json:"activeAt"`
	CreatedAt int     `json:"createdAt"`
}

type Scope struct {
	Type      string `json:"type"`
	Origin    string `json:"origin"`
	CreatedAt int    `json:"createdAt"`
	ExpiresAt int    `json:"expiresAt"`
	TeamId    string `json:"teamId"`
}
