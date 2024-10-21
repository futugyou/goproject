package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type ArtifactService service

type CheckArtifactParameter struct {
	Hash             string
	BaseUrlParameter `json:"-"`
}

func (v *ArtifactService) CheckArtifactExists(ctx context.Context, request CheckArtifactParameter) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v8/artifacts/%s", request.Hash),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type QueryArtifactParameter struct {
	Hashes           []string `json:"hashes"`
	BaseUrlParameter `json:"-"`
}

func (v *ArtifactService) QueryArtifact(ctx context.Context, request QueryArtifactParameter) (*QueryArtifactResponse, error) {
	u := &url.URL{
		Path: "/v8/artifacts",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &QueryArtifactResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DownloadArtifactParameter struct {
	Hash             string
	BaseUrlParameter `json:"-"`
}

func (v *ArtifactService) DownloadArtifact(ctx context.Context, request DownloadArtifactParameter) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v8/artifacts/%s", request.Hash),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *ArtifactService) RecordArtifactEvent(ctx context.Context, request ArtifactEventRequest) (*string, error) {
	u := &url.URL{
		Path: "/v8/artifacts/events",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Post(ctx, path, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type GetCachingStatusParameter struct {
	BaseUrlParameter `json:"-"`
}

func (v *ArtifactService) GetCachingStatus(ctx context.Context, request GetCachingStatusParameter) (*CachingStatus, error) {
	u := &url.URL{
		Path: "/v8/artifacts/status",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CachingStatus{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

// TODO: need common http file upload method first.
// type UploadArtifactRequest struct {
// 	Hash             string
// 	BaseUrlParameter `json:"-"`
// 	// Required head: Content-Length
// 	// Allowed heads: x-artifact-client-ci, x-artifact-client-interactive, x-artifact-duration, x-artifact-tag
// 	CustomeHeader map[string]string
// }

// func (v *ArtifactService) UploadArtifact(ctx context.Context, request UploadArtifactRequest) (*UploadArtifactResponse, error) {
// 	if _, ok := request.CustomeHeader["Content-Length"]; !ok {
// 		return nil, fmt.Errorf("UploadArtifact request MUST have Content-Length parameter in CustomeHeader")
// 	}

// 	u := &url.URL{
// 		Path: "/v8/artifacts/status",
// 	}
// 	params := url.Values{}
// 	if request.TeamId != nil {
// 		params.Add("teamId", *request.TeamId)
// 	}
// 	if request.TeamSlug != nil {
// 		params.Add("slug", *request.TeamSlug)
// 	}
// 	u.RawQuery = params.Encode()
// 	path := u.String()

// 	response := &UploadArtifactResponse{}
// 	client := NewClientWithHeader(request.CustomeHeader)
// 	if err := client.http.Put(ctx, path, nil, response); err != nil {
// 		return nil, err
// 	}
// 	return response, nil
// }

type ArtifactEventRequest struct {
	Event            string `json:"event"`
	Hash             string `json:"hash"`
	SessionId        string `json:"sessionId"`
	Source           string `json:"source"`
	Duration         int    `json:"duration"`
	BaseUrlParameter `json:"-"`
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
