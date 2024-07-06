package vercel

import (
	"fmt"
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

func (v *VercelClient) DeleteAuthToken(tokenId string) (*DeleteAuthTokenResponse, error) {
	path := fmt.Sprintf("/v3/user/tokens/%s", tokenId)
	result := &DeleteAuthTokenResponse{}

	err := v.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) LoginWithEmail(email string, tokenName string) (*LoginWithEmailResponse, error) {
	path := "/registration"
	result := &LoginWithEmailResponse{}
	request := struct {
		Email     string `json:"email"`
		TokenName string `json:"tokenName"`
	}{
		Email:     email,
		TokenName: tokenName,
	}
	err := v.http.Post(path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetAuthTokenMetadata(tokenId string) (*TokenInfo, error) {
	path := fmt.Sprintf("/v5/user/tokens/%s", tokenId)
	result := &TokenInfo{}

	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) ListAuthTokens() (*ListAuthTokensResponse, error) {
	path := "/v5/user/tokens"
	result := &ListAuthTokensResponse{}

	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) VerifyAuthToken(token string, email string) (*VerifyAuthTokenResponse, error) {
	path := "/registration/verify"
	queryParams := url.Values{}
	queryParams.Add("token", token)
	queryParams.Add("email", email)
	path += "?" + queryParams.Encode()

	result := &VerifyAuthTokenResponse{}

	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type VerifyAuthTokenResponse struct {
	Email  string       `json:"email"`
	Token  string       `json:"token"`
	TeamId string       `json:"teamId"`
	Error  *VercelError `json:"error"`
}

type ListAuthTokensResponse struct {
	Pagination   Pagination   `json:"pagination"`
	TestingToken TokenInfo    `json:"testingToken"`
	Tokens       []TokenInfo  `json:"tokens"`
	Error        *VercelError `json:"error"`
}

type LoginWithEmailResponse struct {
	SecurityCode string       `json:"securityCode"`
	Token        string       `json:"token"`
	Error        *VercelError `json:"error"`
}

type DeleteAuthTokenResponse struct {
	TokenId string       `json:"tokenId"`
	Error   *VercelError `json:"error"`
}

type CreateAuthTokenResponse struct {
	BearerToken string       `json:"bearerToken"`
	Token       TokenInfo    `json:"token"`
	Error       *VercelError `json:"error"`
}

type TokenInfo struct {
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Type      string       `json:"type"`
	Origin    string       `json:"origin"`
	Scopes    []Scope      `json:"scopes"`
	ExpiresAt int          `json:"expiresAt"`
	ActiveAt  int          `json:"activeAt"`
	CreatedAt int          `json:"createdAt"`
	Error     *VercelError `json:"error"`
}

type Scope struct {
	Type      string `json:"type"`
	Origin    string `json:"origin"`
	CreatedAt int    `json:"createdAt"`
	ExpiresAt int    `json:"expiresAt"`
	TeamId    string `json:"teamId"`
}
