package lib

import "fmt"

const finetunesPath string = "fine-tunes"
const listFinetunesPath string = "fine-tunes"
const retrieveFintunePath string = "fine-tunes/%s"
const cancelFinetunesPath string = "fine-tunes/%s/cancel"

type CreateFinetuneRequest struct {
	TrainingFile                 string      `json:"training_file"`
	ValidationFile               string      `json:"validation_file,omitempty"`
	Model                        string      `json:"model,omitempty"`
	N_epochs                     int32       `json:"n_epochs,omitempty"`
	BatchSize                    int32       `json:"batch_size,omitempty"`
	LearningRateMultiplier       float32     `json:"learning_rate_multiplier,omitempty"`
	PromptLossWeight             float32     `json:"prompt_loss_weight,omitempty"`
	ComputeClassificationMetrics bool        `json:"compute_classification_metrics,omitempty"`
	ClassificationNClasses       int32       `json:"classification_n_classes,omitempty"`
	ClassificationPositiveClass  string      `json:"classification_positive_class,omitempty"`
	ClassificationBetas          interface{} `json:"classification_betas,omitempty"`
	Suffix                       string      `json:"suffix,omitempty"`
}

type CreateFinetuneResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	FinetuneMoel
}

type FinetuneMoel struct {
	Object          string       `json:"object,omitempty"`
	ID              string       `json:"id,omitempty"`
	Hyperparams     *Hyperparams `json:"hyperparams,omitempty"`
	OrganizationID  string       `json:"organization_id,omitempty"`
	Model           string       `json:"model,omitempty"`
	TrainingFiles   []FileModel  `json:"training_files,omitempty"`
	ValidationFiles []FileModel  `json:"validation_files,omitempty"`
	ResultFiles     []FileModel  `json:"result_files,omitempty"`
	CreatedAt       int          `json:"created_at,omitempty"`
	UpdatedAt       int          `json:"updated_at,omitempty"`
	Status          string       `json:"status,omitempty"`
	FineTunedModel  string       `json:"fine_tuned_model,omitempty"`
	Events          []Events     `json:"events,omitempty"`
}

type Hyperparams struct {
	NEpochs                int32   `json:"n_epochs"`
	BatchSize              int32   `json:"batch_size"`
	PromptLossWeight       float32 `json:"prompt_loss_weight"`
	LearningRateMultiplier float32 `json:"learning_rate_multiplier"`
}

type Events struct {
	Object    string `json:"object"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	CreatedAt int32  `json:"created_at"`
}

type CancelFinetuneResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	FinetuneMoel
}

type ListFinetuneResponse struct {
	Error  *OpenaiError   `json:"error,omitempty"`
	Object string         `json:"object,omitempty"`
	Data   []FinetuneMoel `json:"data,omitempty"`
}

type RetrieveFinetuneResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	FinetuneMoel
}

func (client *openaiClient) CreateFinetune(request CreateFinetuneRequest) *CreateFinetuneResponse {
	result := &CreateFinetuneResponse{}
	client.Post(finetunesPath, request, result)
	return result
}

func (client *openaiClient) CancelFinetune(fine_tune_id string) *CancelFinetuneResponse {
	result := &CancelFinetuneResponse{}
	client.Post(fmt.Sprintf(cancelFinetunesPath, fine_tune_id), nil, result)
	return result
}

func (client *openaiClient) ListFinetune() *ListFinetuneResponse {
	result := &ListFinetuneResponse{}
	client.Get(listFinetunesPath, result)
	return result
}

func (client *openaiClient) RetrieveFinetune(fine_tune_id string) *RetrieveFinetuneResponse {
	result := &RetrieveFinetuneResponse{}
	client.Get(fmt.Sprintf(retrieveFintunePath, fine_tune_id), result)
	return result
}
