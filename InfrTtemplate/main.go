package main

import (
	_ "github.com/joho/godotenv/autoload"
	// "context"
	// "fmt"
	// "os"
	// "github.com/futugyou/InfrTtemplate/core"
	// "github.com/futugyou/InfrTtemplate/entity"
	// "github.com/futugyou/InfrTtemplate/mongorepo"
	// "github.com/futugyou/InfrTtemplate/repository"
)

//go:generate go install github.com/joho/godotenv/cmd/godotenv@latest
//go:generate godotenv -f ./.env go run ../tour/main.go mongo generate
func main() {
	// Test()
	// s := NewBaseDataService()
	// list := s.Paging()
	// for _, v := range list {
	// 	fmt.Println(v.Symbol)
	// }
}

// type BaseDataService struct {
// 	Repository repository.IBaseDatasRepository
// }

// func NewBaseDataService() *BaseDataService {
// 	config := mongorepo.DBConfig{
// 		DBName:        os.Getenv("db_name"),
// 		ConnectString: os.Getenv("mongodb_url"),
// 	}

// 	return &BaseDataService{
// 		Repository: mongorepo.NewBaseDatasRepository(config),
// 	}
// }

// func (s *BaseDataService) Paging() []entity.BaseDatasEntity {
// 	list, err := s.Repository.Paging(context.Background(), core.Paging{Page: 2, Limit: 2}, []core.DataFilterItem{})
// 	if err != nil {
// 		fmt.Println(err)
// 		return []entity.BaseDatasEntity{}
// 	}
// 	return list
// }
