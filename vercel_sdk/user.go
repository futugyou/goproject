package vercel

import (
	"context"
	"net/url"
)

type UserService service

func (v *UserService) GetUser(ctx context.Context) (*AuthUser, error) {
	path := "/v2/user"
	result := &AuthUser{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *UserService) ListUserEvent(ctx context.Context, slug string, teamId string, limit string, since string, until string,
	types string, userId string, withPayload string) ([]UserEvent, error) {
	path := "/v3/events"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(limit) > 0 {
		queryParams.Add("limit", limit)
	}
	if len(since) > 0 {
		queryParams.Add("since", since)
	}
	if len(until) > 0 {
		queryParams.Add("until", until)
	}
	if len(types) > 0 {
		queryParams.Add("types", types)
	}
	if len(userId) > 0 {
		queryParams.Add("userId", userId)
	}
	if len(withPayload) > 0 {
		queryParams.Add("withPayload", withPayload)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := []UserEvent{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *UserService) DeleteUserAccount(ctx context.Context, req DeleteUserRequest) (*DeleteUserResponse, error) {
	path := "/v1/user"
	result := &DeleteUserResponse{}
	err := v.client.http.DeleteWithBody(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type DeleteUserResponse struct {
	Email   string       `json:"email,omitempty"`
	Id      string       `json:"id,omitempty"`
	Message string       `json:"message,omitempty"`
	Error   *VercelError `json:"error,omitempty"`
}

type DeleteUserRequest struct {
	Reasons []DeleteUserReason `json:"reasons,omitempty"`
}

type DeleteUserReason struct {
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug,omitempty"`
}

type AuthUser struct {
	Id                string       `json:"id,omitempty"`
	Email             string       `json:"email,omitempty"`
	Name              string       `json:"name,omitempty"`
	Username          string       `json:"username,omitempty"`
	Avatar            string       `json:"avatar,omitempty"`
	DefaultTeamId     string       `json:"defaultTeamId,omitempty"`
	Version           string       `json:"version,omitempty"`
	Limited           bool         `json:"limited,omitempty"`
	SoftBlock         interface{}  `json:"softBlock,omitempty"`
	CreatedAt         int          `json:"createdAt,omitempty"`
	Billing           interface{}  `json:"billing,omitempty"`
	ResourceConfig    interface{}  `json:"resourceConfig,omitempty"`
	StagingPrefix     string       `json:"stagingPrefix,omitempty"`
	HasTrialAvailable bool         `json:"hasTrialAvailable,omitempty"`
	Error             *VercelError `json:"error,omitempty"`
}

type UserEvent struct {
	Id       string            `json:"id,omitempty"`
	Text     string            `json:"text,omitempty"`
	Entities []UserEventEntity `json:"entities,omitempty"`
	UserId   string            `json:"userId,omitempty"`
	Error    *VercelError      `json:"error,omitempty"`
}

type UserEventEntity struct {
	Type  string `json:"type,omitempty"`
	Start int    `json:"start,omitempty"`
	End   int    `json:"end,omitempty"`
}
