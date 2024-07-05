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

func (v *VercelClient) QueryArtifact(slug string, teamId string, hashes []string) (*QueryArtifactResponse, error) {
	path := "/v8/artifacts"
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
	result := &QueryArtifactResponse{}
	request := struct {
		Hashes []string `json:"hashes"`
	}{
		Hashes: hashes,
	}
	err := v.http.Post(path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) DownloadArtifact(slug string, teamId string, hash string) (*string, error) {
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

func (v *VercelClient) RecordArtifactEvent(teamId string, hash string, request ArtifactEventRequest) (*string, error) {
	path := "/v8/artifacts/events"
	queryParams := url.Values{}

	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := ""
	err := v.http.Post(path, request, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *VercelClient) GetCachingStatus(slug string, teamId string) (*CachingStatus, error) {
	path := "/v8/artifacts/status"
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
	result := &CachingStatus{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type ArtifactEventRequest struct {
	Event     string `json:"event"`
	Hash      string `json:"hash"`
	SessionId string `json:"sessionId"`
	Source    string `json:"source"`
	Duration  int    `json:"duration"`
}

type QueryArtifactResponse struct {
	Size           int          `json:"size"`
	Tag            string       `json:"tag"`
	TaskDurationMs int          `json:"taskDurationMs"`
	Error          *VercelError `json:"error"`
}

type CachingStatus struct {
	Status string       `json:"status"`
	Error  *VercelError `json:"error"`
}
