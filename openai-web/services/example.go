package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/exp/slices"
)

type ExampleModel struct {
	Key              string   `json:"key,omitempty"`
	Title            string   `json:"title,omitempty"`
	SubTitle         string   `json:"subTitle,omitempty"`
	Model            string   `json:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	MaxTokens        int32    `json:"max_tokens,omitempty"`
	Top_p            float32  `json:"top_p,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Description      string   `json:"description,omitempty"`
	SampleResponse   string   `json:"sample_response,omitempty"`
}

type ExampleService struct {
}

const GetAllExamplesKey string = "GetAllExamplesKey"

func (s *ExampleService) GetExampleSettings() []ExampleModel {
	result := make([]ExampleModel, 0)
	// get data from redis,
	// it is not necessary at the moment, but examples.json data will migrate to db in the future
	rmap, _ := Rbd.HGetAll(ctx, GetAllExamplesKey).Result()

	if len(rmap) > 0 {
		for _, r := range rmap {
			example := ExampleModel{}
			json.Unmarshal([]byte(r), &example)
			result = append(result, example)
		}

		return result
	}

	// get data from file
	var examples []byte
	var err error

	if examples, err = os.ReadFile("./examples/examples.json"); err != nil {
		logs.Error(err)
		return result
	}

	if err = json.Unmarshal(examples, &result); err != nil {
		logs.Error(err)
		return result
	}

	examplesCache := make(map[string]interface{})
	for _, example := range result {
		examplestring, _ := json.Marshal(example)
		examplesCache[example.Key] = examplestring
	}

	count, err := Rbd.HSet(ctx, GetAllExamplesKey, examplesCache).Result()
	if err != nil {
		logs.Error(err)
	} else {
		logs.Info(fmt.Sprintf("example data count: %d", count))
	}

	return result
}

func (s *ExampleService) CreateCustomExample(model ExampleModel) []ExampleModel {
	// system examples
	result := s.GetExampleSettings()

	var examples []byte
	var err error

	customExamples := make([]ExampleModel, 0)
	if examples, err = os.ReadFile("./examples/custom.json"); err != nil {
		logs.Error(err)
		return result
	}

	if len(examples) > 0 {
		if err = json.Unmarshal(examples, &customExamples); err != nil {
			logs.Error(err)
			return result
		}
	}

	idx := slices.IndexFunc(customExamples, func(c ExampleModel) bool { return c.Key == model.Key })
	if idx >= 0 {
		return append(result, customExamples...)
	}

	result = append(result, customExamples...)
	customefile, err := os.OpenFile("./examples/custom.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logs.Error(err)
		return result
	}

	defer customefile.Close()

	example, err := json.Marshal(append(customExamples, model))
	if err != nil {
		logs.Error(err)
		return result
	}

	_, err = customefile.Write(example)
	if err != nil {
		logs.Error(err)
		return result
	}
	return append(result, model)
}
