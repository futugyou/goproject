package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type SecurityService service

type ReadFirewallConfigParameter struct {
	ProjectId        string `json:"-"`
	BaseUrlParameter `json:"-"`
}

type ReadFirewallConfig struct {
	FirewallEnabled bool           `json:"firewallEnabled"`
	OwnerId         string         `json:"ownerId"`
	ProjectKey      string         `json:"projectKey"`
	UpdatedAt       string         `json:"updatedAt"`
	Version         int            `json:"version"`
	CRS             map[string]Gen `json:"crs"`
	IPS             []IP           `json:"ips"`
	ManagedRules    ManagedRules   `json:"managedRules"`
	Rules           []Rule         `json:"rules"`
	Changes         any            `json:"changes"`
	Error           *VercelError   `json:"error,omitempty"`
}

func (v *SecurityService) ReadFirewallConfig(ctx context.Context, request ReadFirewallConfigParameter) (*ReadFirewallConfig, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/security/firewall/config/%s", request.ProjectId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ReadFirewallConfig{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type PutFirewallConfigRequest struct {
	ProjectId        string `json:"-"`
	BaseUrlParameter `json:"-"`
	FirewallEnabled  bool           `json:"firewallEnabled"`
	CRS              map[string]Gen `json:"crs"`
	IPS              []IP           `json:"ips"`
	ManagedRules     ManagedRules   `json:"managedRules"`
	Rules            []Rule         `json:"rules"`
}

type PutFirewallConfigResponse struct {
	Active ReadFirewallConfig `json:"active"`
	Error  *VercelError       `json:"error,omitempty"`
}

func (v *SecurityService) PutFirewallConfig(ctx context.Context, request PutFirewallConfigRequest) (*PutFirewallConfigResponse, error) {
	u := &url.URL{
		Path: "/v1/security/firewall/config",
	}
	params := request.GetUrlValues()
	params.Add("projectId", request.ProjectId)
	u.RawQuery = params.Encode()
	path := u.String()

	response := &PutFirewallConfigResponse{}
	if err := v.client.http.Put(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateAttackModeRequest struct {
	BaseUrlParameter      `json:"-"`
	AttackModeEnabled     bool   `json:"attackModeEnabled"`
	ProjectId             string `json:"projectId"`
	AttackModeActiveUntil int    `json:"attackModeActiveUntil"`
}

type UpdateAttackModeResponse struct {
	AttackModeEnabled   bool         `json:"attackModeEnabled"`
	AttackModeUpdatedAt int          `json:"attackModeUpdatedAt"`
	Error               *VercelError `json:"error,omitempty"`
}

func (v *SecurityService) UpdateAttackMode(ctx context.Context, request UpdateAttackModeRequest) (*UpdateAttackModeResponse, error) {
	u := &url.URL{
		Path: "/v1/security/attack-mode",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &UpdateAttackModeResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateFirewallConfigRequest struct {
	ProjectId        string `json:"-"`
	BaseUrlParameter `json:"-"`
	Action           string      `json:"action"`
	Id               string      `json:"id"`
	Value            interface{} `json:"value"`
}

// https://vercel.com/docs/rest-api/endpoints/security#update-firewall-configuration
// Can not found response type in doc
func (v *SecurityService) UpdateFirewallConfig(ctx context.Context, request UpdateFirewallConfigRequest) (*string, error) {
	u := &url.URL{
		Path: "/v1/security/firewall/config",
	}
	params := request.GetUrlValues()
	params.Add("projectId", request.ProjectId)
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Patch(ctx, path, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type Gen struct {
	Active bool   `json:"active"`
	Action string `json:"action"`
}

type IP struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Notes    string `json:"notes"`
	Action   string `json:"action"`
}

type ManagedRules struct {
	Owasp              Owasp `json:"owasp"`
	VerifiedBotsBypass Owasp `json:"verifiedBotsBypass"`
}

type Owasp struct {
	Active    bool   `json:"active"`
	UpdatedAt string `json:"updatedAt"`
	UserId    string `json:"userId"`
	Username  string `json:"username"`
}

type Rule struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Active         bool             `json:"active"`
	ConditionGroup []ConditionGroup `json:"conditionGroup"`
	Action         ActionClass      `json:"action"`
}

type ActionClass struct {
	Mitigate Mitigate `json:"mitigate"`
}

type Mitigate struct {
	Action         string    `json:"action"`
	RateLimit      RateLimit `json:"rateLimit"`
	Redirect       Redirect  `json:"redirect"`
	ActionDuration string    `json:"actionDuration"`
}

type ConditionGroup struct {
	Conditions []Condition `json:"conditions"`
}

type RateLimit struct {
	Action string   `json:"action"`
	Algo   string   `json:"algo"`
	Keys   []string `json:"keys"`
	Limit  int      `json:"limit"`
	Window int      `json:"window"`
}

type Redirect struct {
	Location  string `json:"location"`
	Permanent bool   `json:"permanent"`
}

type Condition struct {
	Key   string      `json:"key"`
	Neg   bool        `json:"neg"`
	Op    string      `json:"op"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
