package lib

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"
)

const finetunesPath string = "fine-tunes"
const listFinetunesPath string = "fine-tunes"
const retrieveFintunePath string = "fine-tunes/%s"
const cancelFinetunesPath string = "fine-tunes/%s/cancel"
const listFinetuneEventPath string = "fine-tunes/%s/events"
const listFinetuneEventStreamPath string = "fine-tunes/%s/events?stream=true"
const deleteFinetuneModelPath string = "models/%s"

var supportedFineTunesModel = []string{Ada, Babbage, Curie, Davinci}

type CreateFinetuneRequest struct {
	TrainingFile   string `json:"training_file"`
	ValidationFile string `json:"validation_file,omitempty"`
	// The name of the base model to fine-tune. You can select one of "ada", "babbage", "curie", "davinci",
	// or a fine-tuned model created after 2022-04-21.
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

type ListFinetuneEventResponse struct {
	Error  *OpenaiError `json:"error,omitempty"`
	Object string       `json:"object,omitempty"`
	Data   []Events     `json:"data,omitempty"`
}

type DeleteFinetuneModelResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	Object  string       `json:"object,omitempty"`
	ID      string       `json:"id,omitempty"`
	Deleted bool         `json:"deleted,omitempty"`
}

func (c *openaiClient) CreateFinetune(request CreateFinetuneRequest) *CreateFinetuneResponse {
	result := &CreateFinetuneResponse{}

	if len(request.Model) > 19 {
		l := request.Model[len(request.Model)-19 : len(request.Model)-9]
		modelDate, err := time.Parse("2006-01-02", l)
		if err != nil {
			result.Error = NewError("fine tune model error", nil)
			return result
		}

		baseDate, _ := time.Parse("2006-01-02", "2022-04-21")
		if baseDate.After(modelDate) {
			result.Error = NewError("fine tune model date can not earlier than 2022-04-21.", nil)
			return result
		}

	} else if len(request.Model) > 0 {
		if !slices.Contains(supportedFineTunesModel, request.Model) {
			result.Error = NewError("fine tune model error", nil)
			return result
		}
	}

	c.httpClient.Post(finetunesPath, request, result)
	return result
}

func (c *openaiClient) CancelFinetune(fine_tune_id string) *CancelFinetuneResponse {
	result := &CancelFinetuneResponse{}
	c.httpClient.Post(fmt.Sprintf(cancelFinetunesPath, fine_tune_id), nil, result)
	return result
}

func (c *openaiClient) ListFinetune() *ListFinetuneResponse {
	result := &ListFinetuneResponse{}
	c.httpClient.Get(listFinetunesPath, result)
	return result
}

func (c *openaiClient) RetrieveFinetune(fine_tune_id string) *RetrieveFinetuneResponse {
	result := &RetrieveFinetuneResponse{}
	c.httpClient.Get(fmt.Sprintf(retrieveFintunePath, fine_tune_id), result)
	return result
}

func (c *openaiClient) ListFinetuneEvents(fine_tune_id string) *ListFinetuneEventResponse {
	result := &ListFinetuneEventResponse{}
	c.httpClient.Get(fmt.Sprintf(listFinetuneEventPath, fine_tune_id), result)
	return result
}

func (c *openaiClient) DeleteFinetuneMdel(model string) *DeleteFinetuneModelResponse {
	result := &DeleteFinetuneModelResponse{}
	c.httpClient.Delete(fmt.Sprintf(deleteFinetuneModelPath, model), result)
	return result
}

func (c *openaiClient) ListFinetuneEventsStream(fine_tune_id string) *ListFinetuneEventResponse {
	result := &ListFinetuneEventResponse{
		Object: "list",
	}

	c.httpClient.GetStream(fmt.Sprintf(listFinetuneEventStreamPath, fine_tune_id))

	defer c.httpClient.Close()

	for {
		if c.httpClient.StreamEnd {
			break
		}

		event := &Events{}
		c.httpClient.ReadStream(event)
		if !c.httpClient.StreamEnd {
			result.Data = append(result.Data, *event)
		}
	}

	return result
}
