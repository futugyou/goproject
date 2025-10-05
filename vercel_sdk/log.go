package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type LogService service

type GetDeploymentLogsRequest struct {
	ProjectId        string `json:"projectId,omitempty"`
	DeploymentId     string `json:"deploymentId,omitempty"`
	BaseUrlParameter `json:"-"`
}

func (v *LogService) GetDeploymentLogs(ctx context.Context, request GetDeploymentLogsRequest) (*LogInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/deployments/%s/runtime-logs", request.ProjectId, request.DeploymentId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &LogInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type LogInfo struct {
	Level              string       `json:"level"`
	Message            string       `json:"message"`
	RowID              string       `json:"rowId"`
	Source             string       `json:"source"`
	TimestampInMS      int64        `json:"timestampInMs"`
	Domain             string       `json:"domain"`
	MessageTruncated   bool         `json:"messageTruncated"`
	RequestMethod      string       `json:"requestMethod"`
	RequestPath        string       `json:"requestPath"`
	ResponseStatusCode int64        `json:"responseStatusCode"`
	Error              *VercelError `json:"error,omitempty"`
}
