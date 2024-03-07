package dynamo2struct

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"

	"github/go-project/tour/internal/common"
)

type Manager struct {
	DB               *dynamodb.Client
	EntityFolder     string
	RepoFolder       string
	Template         *common.Template
	BasePackageName  string
	CoreFoler        string
	DynamoRepoFolder string
}

func NewManager(db *dynamodb.Client, entityFolder string, repoFolder string, pkgName string, coreFoler string, dynamoRepoFolder string) *Manager {
	return &Manager{
		DB:               db,
		EntityFolder:     entityFolder,
		RepoFolder:       repoFolder,
		Template:         common.NewDefaultTemplate(""),
		BasePackageName:  pkgName,
		CoreFoler:        coreFoler,
		DynamoRepoFolder: dynamoRepoFolder,
	}
}

func (m *Manager) Generator() {
	tables, err := m.getAllTable()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(tables)
}

func (m *Manager) getAllTable() ([]string, error) {
	input := dynamodb.ListTablesInput{}

	output, err := m.DB.ListTables(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return output.TableNames, nil
}
