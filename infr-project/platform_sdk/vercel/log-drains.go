package vercel

import (
	"net/url"
)

func (v *VercelClient) CreateLogDrain(slug string, teamId string, req CreateLogDrainRequest) (*LogDrainInfo, error) {
	path := "/v1/log-drains"
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
	result := &LogDrainInfo{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateLogDrainRequest struct {
	DeliveryFormat string      `json:"deliveryFormat"`
	Sources        []string    `json:"sources"`
	Url            string      `json:"url"`
	Environments   []string    `json:"environments"`
	Headers        interface{} `json:"headers"`
	ProjectIds     []string    `json:"projectIds"`
	SamplingRate   int         `json:"samplingRate"`
	Secret         string      `json:"secret"`
}

type LogDrainInfo struct {
	ClientId            string       `json:"clientId"`
	Compression         string       `json:"compression"`
	ConfigurationId     string       `json:"configurationId"`
	CreatedFrom         string       `json:"createdFrom"`
	DeliveryFormat      string       `json:"deliveryFormat"`
	DisabledBy          string       `json:"disabledBy"`
	CreatedAt           int          `json:"createdAt"`
	DeletedAt           int          `json:"deletedAt"`
	DisabledAt          int          `json:"disabledAt"`
	DisabledReason      string       `json:"disabledReason"`
	Environments        []string     `json:"environments"`
	FirstErrorTimestamp int          `json:"firstErrorTimestamp"`
	Headers             interface{}  `json:"headers"`
	Id                  string       `json:"id"`
	Name                string       `json:"name"`
	OwnerId             string       `json:"ownerId"`
	ProjectIds          []string     `json:"projectIds"`
	SamplingRate        int          `json:"samplingRate"`
	Secret              string       `json:"secret"`
	Sources             []string     `json:"sources"`
	Status              string       `json:"status"`
	TeamId              string       `json:"teamId"`
	Url                 string       `json:"url"`
	UpdatedAt           int          `json:"updatedAt"`
	Error               *VercelError `json:"error,omitempty"`
}
