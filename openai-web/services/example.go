package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExampleModel struct {
	Key              string   `json:"key,omitempty" bson:"key,omitempty"`
	Title            string   `json:"title,omitempty" bson:"title,omitempty"`
	SubTitle         string   `json:"subTitle,omitempty" bson:"subTitle,omitempty"`
	Model            string   `json:"model,omitempty" bson:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty" bson:"prompt,omitempty"`
	Temperature      float32  `json:"temperature,omitempty" bson:"temperature,omitempty"`
	MaxTokens        int32    `json:"max_tokens,omitempty" bson:"max_tokens,omitempty"`
	Top_p            float32  `json:"top_p,omitempty" bson:"top_p,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty" bson:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty" bson:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty" bson:"stop,omitempty"`
	Tags             []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Description      string   `json:"description,omitempty" bson:"description,omitempty"`
	SampleResponse   string   `json:"sample_response,omitempty" bson:"sample_response,omitempty"`
}

type ExampleService struct {
}

const GetAllExamplesKey string = "GetAllExamplesKey"

func (s *ExampleService) GetSystemExamples() []ExampleModel {
	result := make([]ExampleModel, 0)
	// get data from redis,
	// it is not necessary at the moment, but examples.json data will migrate to db in the future
	rmap, e := Rbd.HGetAll(ctx, GetAllExamplesKey).Result()

	if e != nil {
		fmt.Println(e)
	}
	if len(rmap) > 0 {
		for _, r := range rmap {
			example := ExampleModel{}
			json.Unmarshal([]byte(r), &example)
			result = append(result, example)
		}

		return result
	}

	result = s.getExamples("examples")

	examplesCache := make(map[string]interface{})
	for _, example := range result {
		examplestring, _ := json.Marshal(example)
		examplesCache[example.Key] = examplestring
	}

	count, err := Rbd.HSet(ctx, GetAllExamplesKey, examplesCache).Result()
	if err != nil {
		logs.Error(err)
	} else {
		_, err := Rbd.Expire(ctx, GetAllExamplesKey, time.Hour).Result()
		if err != nil {
			logs.Error(err)
		} else {
			logs.Info(fmt.Sprintf("example data count: %d", count))
		}
	}

	return result
}

func (s *ExampleService) CreateSystemExample(model ExampleModel) {
	s.createExample(model, "examples")
}

func (s *ExampleService) GetCustomExamples() []ExampleModel {
	return s.getExamples("examples_custome")
}

func (s *ExampleService) CreateCustomExample(model ExampleModel) {
	s.createExample(model, "examples_custome")
}

func (s *ExampleService) createExample(model ExampleModel, tableName string) {
	uri := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(db_name).Collection(tableName)
	var example ExampleModel
	upsert := true
	option := options.FindOneAndReplaceOptions{
		Upsert: &upsert,
	}

	err = coll.FindOneAndReplace(context.TODO(), bson.D{{Key: "key", Value: model.Key}}, model, &option).Decode(&example)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}
	}
}

func (s *ExampleService) getExamples(tableName string) []ExampleModel {
	result := make([]ExampleModel, 0)

	uri := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(db_name).Collection(tableName)
	filter := bson.D{}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	// end find

	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result
}

func (s *ExampleService) InitExamples() {
	result := make([]ExampleModel, 0)

	// get data from file
	var examples []byte
	var err error

	if examples, err = os.ReadFile("./examples/examples.json"); err != nil {
		logs.Error(err)
		return
	}

	if err = json.Unmarshal(examples, &result); err != nil {
		logs.Error(err)
		return
	}

	uri := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(db_name).Collection("examples_raw")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println(err)
		return
	}

	if cursor.TryNext(context.TODO()) {
		return
	}

	newResults := make([]interface{}, len(result))
	for i, v := range result {
		newResults[i] = v
	}

	_, err = coll.InsertMany(context.TODO(), newResults)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *ExampleService) deleteAllExample(tableName string) {
	uri := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(db_name).Collection(tableName)
	filter := bson.D{}
	if _, err = coll.DeleteMany(context.TODO(), filter); err != nil {
		panic(err)
	}
}

func (s *ExampleService) Reset() {

}
