package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type AuthService service

type CreateAuthTokenRequest struct {
	Name             string `json:"name"`
	ExpiresAt        int    `json:"expiresAt"`
	BaseUrlParameter `json:"-"`
}

func (v *AuthService) CreateAuthToken(ctx context.Context, request CreateAuthTokenRequest) (*CreateAuthTokenResponse, error) {
	u := &url.URL{
		Path: "/v3/user/tokens",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CreateAuthTokenResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (v *AuthService) DeleteAuthToken(ctx context.Context, tokenId string) (*DeleteAuthTokenResponse, error) {
	path := fmt.Sprintf("/v3/user/tokens/%s", tokenId)
	response := &DeleteAuthTokenResponse{}

	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

// Deprecated
func (v *AuthService) LoginWithEmail(ctx context.Context, email string, tokenName string) (*LoginWithEmailResponse, error) {
	path := "/registration"
	response := &LoginWithEmailResponse{}
	request := struct {
		Email     string `json:"email"`
		TokenName string `json:"tokenName"`
	}{
		Email:     email,
		TokenName: tokenName,
	}

	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (v *AuthService) GetAuthTokenMetadata(ctx context.Context, tokenId string) (*TokenInfo, error) {
	path := fmt.Sprintf("/v5/user/tokens/%s", tokenId)
	response := &TokenInfo{}

	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (v *AuthService) ListAuthTokens(ctx context.Context) (*ListAuthTokensResponse, error) {
	path := "/v5/user/tokens"
	response := &ListAuthTokensResponse{}

	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

// Deprecated
func (v *AuthService) VerifyAuthToken(ctx context.Context, token string, email string) (*VerifyAuthTokenResponse, error) {
	path := "/registration/verify"
	queryParams := url.Values{}
	queryParams.Add("token", token)
	queryParams.Add("email", email)
	path += "?" + queryParams.Encode()

	response := &VerifyAuthTokenResponse{}

	err := v.client.http.Get(ctx, path, response)

	if err != nil {
		return nil, err
	}
	return response, nil
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
