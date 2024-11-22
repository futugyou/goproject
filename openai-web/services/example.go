package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/redis/go-redis/v9"
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
	db      *mongo.Database
	redisDb *redis.Client
}

func NewExampleService(client *mongo.Client, redisDb *redis.Client) *ExampleService {
	db_name := os.Getenv("db_name")
	if client == nil {
		uri := os.Getenv("mongodb_url")
		var err error
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
	}

	if redisDb == nil {
		client, err := RedisClient(os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}
		redisDb = client
	}

	return &ExampleService{
		db:      client.Database(db_name),
		redisDb: redisDb,
	}
}

const GetAllExamplesKey string = "GetAllExamplesKey"
const ExampleRawTableName string = "examples_raw"
const ExampleTableName string = "examples"
const ExampleCustomeTableName string = "examples_custome"

func (s *ExampleService) GetSystemExamples(ctx context.Context) []ExampleModel {
	result := make([]ExampleModel, 0)
	// get data from redis,
	// it is not necessary at the moment, but examples.json data will migrate to db in the future
	rmap, e := s.redisDb.HGetAll(ctx, GetAllExamplesKey).Result()

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

	result = s.getExamples(ctx, ExampleTableName)

	examplesCache := make(map[string]interface{})
	for _, example := range result {
		examplestring, _ := json.Marshal(example)
		examplesCache[example.Key] = examplestring
	}

	count, err := s.redisDb.HSet(ctx, GetAllExamplesKey, examplesCache).Result()
	if err != nil {
		logs.Error(err)
	} else {
		_, err := s.redisDb.Expire(ctx, GetAllExamplesKey, time.Hour).Result()
		if err != nil {
			logs.Error(err)
		} else {
			logs.Info(fmt.Sprintf("example data count: %d", count))
		}
	}

	return result
}

func (s *ExampleService) CreateSystemExample(ctx context.Context, model ExampleModel) {
	s.createExample(ctx, model, ExampleTableName)
}

func (s *ExampleService) GetCustomExamples(ctx context.Context) []ExampleModel {
	return s.getExamples(ctx, ExampleCustomeTableName)
}

func (s *ExampleService) CreateCustomExample(ctx context.Context, model ExampleModel) {
	s.createExample(ctx, model, ExampleCustomeTableName)
}

func (s *ExampleService) createExample(ctx context.Context, model ExampleModel, tableName string) {
	coll := s.db.Collection(tableName)
	var example ExampleModel
	upsert := true
	option := options.FindOneAndReplaceOptions{
		Upsert: &upsert,
	}

	err := coll.FindOneAndReplace(ctx, bson.D{{Key: "key", Value: model.Key}}, model, &option).Decode(&example)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}
	}
}

func (s *ExampleService) getExamples(ctx context.Context, tableName string) []ExampleModel {
	result := make([]ExampleModel, 0)
	coll := s.db.Collection(tableName)
	filter := bson.D{}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(ctx, &result); err != nil {
		panic(err)
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result
}

func (s *ExampleService) InitExamples(ctx context.Context) {
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

	s.insertManyExample(ctx, ExampleRawTableName, result)
}

func (s *ExampleService) deleteAllExample(ctx context.Context, tableName string) {
	coll := s.db.Collection(tableName)
	filter := bson.D{}
	if _, err := coll.DeleteMany(ctx, filter); err != nil {
		fmt.Println(err)
	}
}

func (s *ExampleService) insertManyExample(ctx context.Context, tableName string, datas []ExampleModel) {
	coll := s.db.Collection(tableName)
	newResults := make([]interface{}, len(datas))
	for i, v := range datas {
		newResults[i] = v
	}

	if _, err := coll.InsertMany(ctx, newResults); err != nil {
		fmt.Println(err)
	}
}

func (s *ExampleService) Reset(ctx context.Context) {
	s.deleteAllExample(ctx, ExampleTableName)
	datas := s.getExamples(ctx, ExampleRawTableName)
	s.insertManyExample(ctx, ExampleTableName, datas)
}
