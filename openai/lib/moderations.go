package lib

const createModerationPath string = "moderations"

type CreateModerationRequest struct {
	Input string `json:"input,omitempty"`
	Model string `json:"model,omitempty"`
}

type CreateModerationResponse struct {
	Error   *OpenaiError       `json:"error,omitempty"`
	ID      string             `json:"id,omitempty"`
	Model   string             `json:"model,omitempty"`
	Results []ModerationResult `json:"results,omitempty"`
}

type Categories struct {
	Sexual          bool `json:"sexual"`
	Hate            bool `json:"hate"`
	Violence        bool `json:"violence"`
	SelfHarm        bool `json:"self-harm"`
	SexualMinors    bool `json:"sexual/minors"`
	HateThreatening bool `json:"hate/threatening"`
	ViolenceGraphic bool `json:"violence/graphic"`
}

type CategoryScores struct {
	Sexual          float64 `json:"sexual,omitempty"`
	Hate            float64 `json:"hate,omitempty"`
	Violence        float64 `json:"violence,omitempty"`
	SelfHarm        float64 `json:"self-harm,omitempty"`
	SexualMinors    float64 `json:"sexual/minors,omitempty"`
	HateThreatening float64 `json:"hate/threatening,omitempty"`
	ViolenceGraphic float64 `json:"violence/graphic,omitempty"`
}

type ModerationResult struct {
	Flagged        bool            `json:"flagged"`
	Categories     *Categories     `json:"categories,omitempty"`
	CategoryScores *CategoryScores `json:"category_scores,omitempty"`
}

func (client *openaiClient) CreateModeration(request CreateModerationRequest) *CreateModerationResponse {
	result := &CreateModerationResponse{}
	client.Post(createModerationPath, request, result)
	return result
}
