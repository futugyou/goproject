package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *dynamodb.Client
)

func init() {
	svc = dynamodb.NewFromConfig(awsenv.Cfg)
}

func ListGlobalTables() {
	input := dynamodb.ListGlobalTablesInput{
		RegionName: aws.String("ap-northeast-1"),
	}
	output, err := svc.ListGlobalTables(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, table := range output.GlobalTables {
		fmt.Println(*table.GlobalTableName, table.ReplicationGroup)
	}
}

func ListTables() {
	input := dynamodb.ListTablesInput{}
	output, err := svc.ListTables(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, table := range output.TableNames {
		fmt.Println(table)
		DescribeTable(table)
	}
}

func DescribeTable(tableName string) {
	input := dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	output, err := svc.DescribeTable(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, attr := range output.Table.AttributeDefinitions {
		fmt.Print("\t", *attr.AttributeName)
	}
	fmt.Print("\t",
		output.Table.BillingModeSummary.BillingMode.Values(), "\t",
		output.Table.CreationDateTime, "\t",
		*output.Table.ItemCount, "\t",
		// *output.Table.TableArn, "\t",
		output.Table.TableStatus)
	for _, schema := range output.Table.KeySchema {
		fmt.Print("\t", *schema.AttributeName, "\t", schema.KeyType)
	}
	fmt.Println()
}
