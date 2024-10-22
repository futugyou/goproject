package vercel

import (
	"context"
	"net/url"
)

type UserService service

func (v *UserService) GetUser(ctx context.Context) (*AuthUser, error) {
	path := "/v2/user"
	response := &AuthUser{}
	err := v.client.http.Get(ctx, path, response)

	if err != nil {
		return nil, err
	}
	return response, nil
}

type ListUserEventParameter struct {
	Limit            *string
	Since            *string
	Until            *string
	Types            *string
	UserId           *string
	WithPayload      *string
	BaseUrlParameter `json:"-"`
}

func (v *UserService) ListUserEvent(ctx context.Context, request ListUserEventParameter) ([]UserEvent, error) {
	u := &url.URL{
		Path: "/v3/events",
	}
	params := request.GetUrlValues()
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	if request.Types != nil {
		params.Add("types", *request.Types)
	}
	if request.UserId != nil {
		params.Add("userId", *request.UserId)
	}
	if request.WithPayload != nil {
		params.Add("withPayload", *request.WithPayload)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []UserEvent{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (v *UserService) DeleteUserAccount(ctx context.Context, req DeleteUserRequest) (*DeleteUserResponse, error) {
	path := "/v1/user"
	response := &DeleteUserResponse{}
	if err := v.client.http.DeleteWithBody(ctx, path, req, response); err != nil {
		return nil, err
	}
	return response, nil
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
