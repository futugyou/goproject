package lib

const finetunesPath string = "fine-tunes"

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
	Error           *OpenaiError      `json:"error,omitempty"`
	Object          string            `json:"object,omitempty"`
	ID              string            `json:"id,omitempty"`
	Hyperparams     *Hyperparams      `json:"hyperparams,omitempty"`
	OrganizationID  string            `json:"organization_id,omitempty"`
	Model           string            `json:"model,omitempty"`
	TrainingFiles   []TrainingFiles   `json:"training_files,omitempty"`
	ValidationFiles []ValidationFiles `json:"validation_files,omitempty"`
	ResultFiles     []fileModel       `json:"result_files,omitempty"`
	CreatedAt       int32             `json:"created_at,omitempty"`
	UpdatedAt       int32             `json:"updated_at,omitempty"`
	Status          string            `json:"status,omitempty"`
	FineTunedModel  string            `json:"fine_tuned_model,omitempty"`
	Events          []Events          `json:"events,omitempty"`
}

type Hyperparams struct {
	NEpochs                int32   `json:"n_epochs"`
	BatchSize              int32   `json:"batch_size"`
	PromptLossWeight       float32 `json:"prompt_loss_weight"`
	LearningRateMultiplier float32 `json:"learning_rate_multiplier"`
}

type TrainingFiles struct {
	Object        string      `json:"object"`
	ID            string      `json:"id"`
	Purpose       string      `json:"purpose"`
	Filename      string      `json:"filename"`
	Bytes         int32       `json:"bytes"`
	CreatedAt     int32       `json:"created_at"`
	Status        string      `json:"status"`
	StatusDetails interface{} `json:"status_details"`
}

type ValidationFiles struct {
	Object        string      `json:"object"`
	ID            string      `json:"id"`
	Purpose       string      `json:"purpose"`
	Filename      string      `json:"filename"`
	Bytes         int32       `json:"bytes"`
	CreatedAt     int32       `json:"created_at"`
	Status        string      `json:"status"`
	StatusDetails interface{} `json:"status_details"`
}

type Events struct {
	Object    string `json:"object"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	CreatedAt int32  `json:"created_at"`
}

func (client *openaiClient) CreateFinetune(request CreateFinetuneRequest) *CreateFinetuneResponse {
	result := &CreateFinetuneResponse{}
	client.Post(finetunesPath, request, result)
	return result
}
