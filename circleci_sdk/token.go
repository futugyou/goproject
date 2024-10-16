package circleci

type TokenService service

func (s *TokenService) DeleteOrgLevelClaims(orgID string, claims []string) (*OrgLevelClaim, error) {
	path := "/org/" + orgID + "/oidc-custom-claims"
	for i := 0; i < len(claims); i++ {
		if i == 0 {
			path += "?"
		} else {
			path += "&"
		}
		path += ("claims=" + claims[i])
	}

	result := &OrgLevelClaim{}
	if err := s.client.http.Delete(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) GetOrgLevelClaims(orgID string) (*OrgLevelClaim, error) {
	path := "/org/" + orgID + "/oidc-custom-claims"

	result := &OrgLevelClaim{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) PatchOrgLevelClaims(orgID string, audience []string, ttl string) (*OrgLevelClaim, error) {
	path := "/org/" + orgID + "/oidc-custom-claims"
	request := UpdateClaimRequest{
		Audience: audience,
		TTL:      ttl,
	}
	result := &OrgLevelClaim{}
	if err := s.client.http.Patch(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) DeleteProjectLevelClaims(orgID string, projectId string, claims []string) (*OrgLevelClaim, error) {
	path := "/org/" + orgID + "/project/" + projectId + "/oidc-custom-claims"
	for i := 0; i < len(claims); i++ {
		if i == 0 {
			path += "?"
		} else {
			path += "&"
		}
		path += ("claims=" + claims[i])
	}

	result := &OrgLevelClaim{}
	if err := s.client.http.Delete(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) GetProjectLevelClaims(orgID string, projectId string) (*OrgLevelClaim, error) {
	path := "/org/" + orgID + "/project/" + projectId + "/oidc-custom-claims"

	result := &OrgLevelClaim{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TokenService) PatchProjectLevelClaims(orgID string, projectId string, audience []string, ttl string) (*OrgLevelClaim, error) {
	path := "/org/" + orgID + "/project/" + projectId + "/oidc-custom-claims"
	request := UpdateClaimRequest{
		Audience: audience,
		TTL:      ttl,
	}
	result := &OrgLevelClaim{}
	if err := s.client.http.Patch(path, request, result); err != nil {
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
