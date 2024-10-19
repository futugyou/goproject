package circleci

import (
	"context"
	"fmt"
)

type TokenService service

func (s *TokenService) DeleteOrgLevelClaims(ctx context.Context, orgID string, claims []string) (*OrgLevelClaim, error) {
	path := fmt.Sprintf("/org/%s/oidc-custom-claims", orgID)
	for i := 0; i < len(claims); i++ {
		if i == 0 {
			path += "?"
		} else {
			path += "&"
		}
		path += ("claims=" + claims[i])
	}

	result := &OrgLevelClaim{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) GetOrgLevelClaims(ctx context.Context, orgID string) (*OrgLevelClaim, error) {
	path := fmt.Sprintf("/org/%s/oidc-custom-claims", orgID)

	result := &OrgLevelClaim{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) PatchOrgLevelClaims(ctx context.Context, orgID string, audience []string, ttl string) (*OrgLevelClaim, error) {
	path := fmt.Sprintf("/org/%s/oidc-custom-claims", orgID)
	request := UpdateClaimRequest{
		Audience: audience,
		TTL:      ttl,
	}
	result := &OrgLevelClaim{}
	if err := s.client.http.Patch(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) DeleteProjectLevelClaims(ctx context.Context, orgID string, projectId string, claims []string) (*OrgLevelClaim, error) {
	path := fmt.Sprintf("/org/%s/project/%s/oidc-custom-claims", orgID, projectId)
	for i := 0; i < len(claims); i++ {
		if i == 0 {
			path += "?"
		} else {
			path += "&"
		}
		path += ("claims=" + claims[i])
	}

	result := &OrgLevelClaim{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) GetProjectLevelClaims(ctx context.Context, orgID string, projectId string) (*OrgLevelClaim, error) {
	path := fmt.Sprintf("/org/%s/project/%s/oidc-custom-claims", orgID, projectId)
	result := &OrgLevelClaim{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) PatchProjectLevelClaims(ctx context.Context, orgID string, projectId string, audience []string, ttl string) (*OrgLevelClaim, error) {
	path := fmt.Sprintf("/org/%s/project/%s/oidc-custom-claims", orgID, projectId)
	request := UpdateClaimRequest{
		Audience: audience,
		TTL:      ttl,
	}
	result := &OrgLevelClaim{}
	if err := s.client.http.Patch(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

type UpdateClaimRequest struct {
	Audience []string `json:"audience"`
	TTL      string   `json:"ttl"`
}

type OrgLevelClaim struct {
	Audience          []string `json:"audience"`
	AudienceUpdatedAt string   `json:"audience_updated_at"`
	OrgID             string   `json:"org_id"`
	ProjectID         string   `json:"project_id"`
	TTL               string   `json:"ttl"`
	TTLUpdatedAt      string   `json:"ttl_updated_at"`
	Message           *string  `json:"message,omitempty"`
}
