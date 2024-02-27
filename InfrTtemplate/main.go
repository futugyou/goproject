package main

import (
	_ "github.com/joho/godotenv/autoload"
	// "context"
	// "fmt"
	// "os"
	// "github.com/futugyou/InfrTtemplate/core"
	// "github.com/futugyou/InfrTtemplate/mongorepo"
)

//go:generate go install github.com/joho/godotenv/cmd/godotenv@latest
//go:generate godotenv -f ./.env go run ../tour/main.go  mongo generate
func main() {
	// config := mongorepo.DBConfig{
	// 	DBName:        os.Getenv("db_name"),
	// 	ConnectString: os.Getenv("mongodb_url"),
	// }

	// repository := mongorepo.NewBaseDatasRepository(config)
	// list, err := repository.Paging(context.Background(), core.Paging{Page: 1, Limit: 100}, []core.DataFilterItem{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for _, v := range list {
	// 	fmt.Println(v.Symbol)
	// }
}
