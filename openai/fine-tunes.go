package openai

import (
	"context"
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

var supportedFineTunesModel = []string{
	GPT3_davinci,
	GPT3_curie,
	GPT3_babbage,
	GPT3_ada,
}

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

type FinetuneService service

func (c *FinetuneService) CreateFinetune(ctx context.Context, request CreateFinetuneRequest) *CreateFinetuneResponse {
	result := &CreateFinetuneResponse{}

	if len(request.Model) > 19 {
		l := request.Model[len(request.Model)-19 : len(request.Model)-9]
		modelDate, err := time.Parse("2006-01-02", l)
		if err != nil {
			result.Error = messageError("fine tune model format error, plaese check your model.")
			return result
		}

		baseDate, _ := time.Parse("2006-01-02", "2022-04-21")
		if baseDate.After(modelDate) {
			result.Error = messageError(fmt.Sprintf("fine tune model date can not earlier than 2022-04-21, current is %s", modelDate.Format("2006-01-02")))
			return result
		}

	} else if len(request.Model) > 0 {
		if !slices.Contains(supportedFineTunesModel, request.Model) {
			result.Error = unsupportedTypeError("Model", request.Model, supportedFineTunesModel)
			return result
		}
	}

	c.client.httpClient.Post(ctx, finetunesPath, request, result)
	return result
}

func (c *FinetuneService) CancelFinetune(ctx context.Context, fine_tune_id string) *CancelFinetuneResponse {
	result := &CancelFinetuneResponse{}
	c.client.httpClient.Post(ctx, fmt.Sprintf(cancelFinetunesPath, fine_tune_id), nil, result)
	return result
}

func (c *FinetuneService) ListFinetune(ctx context.Context) *ListFinetuneResponse {
	result := &ListFinetuneResponse{}
	c.client.httpClient.Get(ctx, listFinetunesPath, result)
	return result
}

func (c *FinetuneService) RetrieveFinetune(ctx context.Context, fine_tune_id string) *RetrieveFinetuneResponse {
	result := &RetrieveFinetuneResponse{}
	c.client.httpClient.Get(ctx, fmt.Sprintf(retrieveFintunePath, fine_tune_id), result)
	return result
}

func (c *FinetuneService) ListFinetuneEvents(ctx context.Context, fine_tune_id string) *ListFinetuneEventResponse {
	result := &ListFinetuneEventResponse{}
	c.client.httpClient.Get(ctx, fmt.Sprintf(listFinetuneEventPath, fine_tune_id), result)
	return result
}

func (c *FinetuneService) DeleteFinetuneMdel(ctx context.Context, model string) *DeleteFinetuneModelResponse {
	result := &DeleteFinetuneModelResponse{}
	c.client.httpClient.Delete(ctx, fmt.Sprintf(deleteFinetuneModelPath, model), result)
	return result
}

// you can read stream in this way.
//
// stream, err := openai.ListFinetuneEventsStream(fine_tune_id)
//
//	if err != nil {
//		doSomething()
//	}
//
// defer stream.Close()
//
// result := &ListFinetuneEventResponse{Object: "list"}
//
//	for {
//		if !stream.CanReadStream() {
//			break
//		}
//		event := &Events{}
//		if err = stream.ReadStream(event); err != nil {
//			doSomething()
//		} else {
//			result.Data = append(result.Data, *event)
//		}
//	}
func (c *FinetuneService) ListFinetuneEventsStream(ctx context.Context, fine_tune_id string) (*StreamResponse, *OpenaiError) {
	return c.client.httpClient.GetStream(ctx, fmt.Sprintf(listFinetuneEventStreamPath, fine_tune_id))
}
