package dynamo2struct

import (
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

type TableInfo struct {
	Name    string
	ColInfo map[string]string
}

func (m *Manager) Generator() {
	tables, err := m.getAllTable()
	if err != nil {
		log.Println(err)
		return
	}

	tInfos := m.concurrentGetTableColum(tables)
	for _, v := range tInfos {
		fmt.Println(v.Name, v.ColInfo)
	}
}

func (m *Manager) concurrentGetTableColum(tables []string) []TableInfo {
	tInfos := make([]TableInfo, 0)
	var wg sync.WaitGroup
	for _, table := range tables {
		wg.Add(1)
		go func(table string, tInfos *[]TableInfo, wg *sync.WaitGroup) {
			defer wg.Done()

			var cols map[string]string
			var err error
			if cols, err = m.getTableColum(table); err != nil {
				log.Println(err)
				return
			}

			*tInfos = append(*tInfos, TableInfo{
				Name:    table,
				ColInfo: cols,
			})
		}(table, &tInfos, &wg)
	}

	wg.Wait()
	return tInfos
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

func (m *Manager) getTableColum(tableName string) (map[string]string, error) {
	result := make(map[string]string)
	// describeInput := dynamodb.DescribeTableInput{
	// 	TableName: aws.String(tableName),
	// }

	// describeOutput, err := m.DB.DescribeTable(awsenv.EmptyContext, &describeInput)
	// if err != nil {
	// 	return result, err
	// }

	// for _, schema := range describeOutput.Table.KeySchema {
	// 	if schema.AttributeName != nil {
	// 		result[*schema.AttributeName] = "string"
	// 	}
	// }

	input := dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int32(10),
	}

	output, err := m.DB.Scan(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	for _, item := range output.Items {
		for key, value := range item {
			valueType := "inferface{}"
			switch value.(type) {
			case *types.AttributeValueMemberB:
				valueType = "[]byte"
			case *types.AttributeValueMemberM:
				valueType = "map[string]interface{}"
			case *types.AttributeValueMemberBOOL:
				valueType = "bool"
			case *types.AttributeValueMemberN:
				valueType = "string"
			case *types.AttributeValueMemberNS:
				valueType = "[]string"
			case *types.AttributeValueMemberS:
				valueType = "string"
			case *types.AttributeValueMemberSS:
				valueType = "[]string"
			}

			if _, ok := result[key]; !ok {
				result[key] = valueType
			}
		}
	}

	return result, nil
}
