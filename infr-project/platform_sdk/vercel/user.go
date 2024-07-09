package vercel

func (v *VercelClient) GetUser() (*AuthUser, error) {
	path := "/v2/user"
	result := &AuthUser{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
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
