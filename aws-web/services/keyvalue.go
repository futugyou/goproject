package services

import (
	"context"
	"os"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
)

type KeyValueService struct {
	repository repository.IKeyValueRepository
}

func NewKeyValueService() *KeyValueService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &KeyValueService{
		repository: mongorepo.NewKeyValueRepository(config),
	}
}

func (a *KeyValueService) GetAllKeyValues() []model.KeyValue {
	KeyValues := make([]model.KeyValue, 0)
	entities, err := a.repository.GetAll(context.Background())
	if err != nil {
		return KeyValues
	}

	for _, entity := range entities {
		KeyValues = append(KeyValues, model.KeyValue{
			Key:   entity.Id,
			Value: entity.Value,
		})
	}

	return KeyValues
}

func (a *KeyValueService) GetValueByKey(key string) *model.KeyValue {
	keyvalue, err := a.repository.Get(context.Background(), key)
	if err != nil {
		return nil
	}

	return &model.KeyValue{
		Key:   keyvalue.Id,
		Value: keyvalue.Value,
	}
}

func (a *KeyValueService) CreateKeyValue(key string, value string) {
	keyvalue, err := a.repository.Get(context.Background(), key)
	if err != nil {
		keyvalue = &entity.KeyValueEntity{
			Id:    key,
			Value: value,
		}
		a.repository.Insert(context.Background(), *keyvalue)
		return
	} else {
		keyvalue.Value = value
		a.repository.Update(context.Background(), *keyvalue, keyvalue.Id)
	}
}
