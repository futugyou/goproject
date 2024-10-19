package vercel

import (
	"context"
	"fmt"
	"net/url"
)

func (v *ArtifactService) CheckArtifactExists(ctx context.Context, hash string, slug string, teamId string) (*string, error) {
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
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ArtifactService) QueryArtifact(ctx context.Context, slug string, teamId string, hashes []string) (*QueryArtifactResponse, error) {
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
	err := v.client.http.Post(ctx, path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ArtifactService) DownloadArtifact(ctx context.Context, slug string, teamId string, hash string) (*string, error) {
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
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ArtifactService) RecordArtifactEvent(ctx context.Context, teamId string, hash string, request ArtifactEventRequest) (*string, error) {
	path := "/v8/artifacts/events"
	queryParams := url.Values{}

	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := ""
	err := v.client.http.Post(ctx, path, request, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ArtifactService) GetCachingStatus(ctx context.Context, slug string, teamId string) (*CachingStatus, error) {
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
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type ArtifactService service

func (v *ArtifactService) UploadArtifact(ctx context.Context, hash string, slug string, teamId string) (*UploadArtifactResponse, error) {
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
	result := &UploadArtifactResponse{}
	//TODO Header Params Content-Length Required
	err := v.client.http.Put(ctx, path, nil, result)

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

type UploadArtifactResponse struct {
	Urls  []string     `json:"urls"`
	Error *VercelError `json:"error"`
}
