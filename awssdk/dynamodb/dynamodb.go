package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func CreateTable() {
	input := dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("LockId"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("LockId"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName:   aws.String("Terraform-Lock"),
		BillingMode: types.BillingModePayPerRequest,
	}
	output, err := svc.CreateTable(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*output.TableDescription.TableId)
}

func CreateBackup() {
	input := dynamodb.CreateBackupInput{
		BackupName: aws.String("Terraform-Lock-Back"),
		TableName:  aws.String("Terraform-Lock"),
	}
	output, err := svc.CreateBackup(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.BackupDetails.BackupType, output.BackupDetails.BackupStatus, output.BackupDetails.BackupCreationDateTime, output.BackupDetails.BackupExpiryDateTime)
}
