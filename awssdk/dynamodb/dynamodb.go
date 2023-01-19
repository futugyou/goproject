package dynamodb

import (
	"fmt"
	"time"

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

func ListBackups() {
	input := dynamodb.ListBackupsInput{
		TableName: aws.String("Terraform-Lock"),
	}
	output, err := svc.ListBackups(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, backup := range output.BackupSummaries {
		fmt.Println(*backup.BackupName)
		DeleteBackup(*backup.BackupArn)
	}
}

func DeleteBackup(backupArn string) {
	input := dynamodb.DeleteBackupInput{
		BackupArn: aws.String(backupArn),
	}
	output, err := svc.DeleteBackup(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.BackupDescription)
}

func DeleteTable() {
	input := dynamodb.DeleteTableInput{
		TableName: aws.String("Terraform-Lock"),
	}
	output, err := svc.DeleteTable(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.TableDescription)
}

func Scan(tableName string) {
	var attrValue types.AttributeValueMemberS = types.AttributeValueMemberS{Value: "some value"}
	input := dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ConsistentRead:            aws.Bool(false), // default false
		FilterExpression:          aws.String("#P = :val"),
		ExpressionAttributeNames:  map[string]string{"#P": "PK"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":val": (types.AttributeValue)(&attrValue)},
	}

	output, err := svc.Scan(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("count:", output.Count)
	fmt.Println(output.LastEvaluatedKey)
	fmt.Println(output.ConsumedCapacity)
	fmt.Println(output.ScannedCount)
	for _, item := range output.Items {
		for key, value := range item {
			fmt.Println(key, value)
		}
		fmt.Println()
	}
}

func Query(tableName string) {
	var attrValue types.AttributeValueMemberS = types.AttributeValueMemberS{Value: "some value"}
	input := dynamodb.QueryInput{
		TableName:     aws.String(tableName),
		KeyConditions: map[string]types.Condition{"PK": {ComparisonOperator: types.ComparisonOperatorEq, AttributeValueList: []types.AttributeValue{&attrValue}}},
	}
	output, err := svc.Query(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range output.Items {
		for key, value := range item {
			fmt.Println(key, value)
		}
		fmt.Println()
	}
}

func GetItem() {
	tableName := "some value"
	var pk types.AttributeValueMemberS = types.AttributeValueMemberS{Value: "some value"}
	var sk types.AttributeValueMemberS = types.AttributeValueMemberS{Value: "some value"}
	input := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &pk,
			"SK": &sk,
		},
	}
	output, err := svc.GetItem(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ConsumedCapacity:", output.ConsumedCapacity)
	for key, value := range output.Item {
		fmt.Println(key, value)
	}
}

// this mothed will use 'Item' to replace all data
func PutItem() {
	var attrValue types.AttributeValueMemberS = types.AttributeValueMemberS{Value: "some value"}
	var time types.AttributeValueMemberS = types.AttributeValueMemberS{Value: time.Now().Format("2006/01/02 15:04:05")}
	input := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"LockTime": &time,
			"PK":       &attrValue,
		},
		TableName:                 aws.String("some name"),
		ConditionExpression:       aws.String("#P = :val"),
		ExpressionAttributeNames:  map[string]string{"#P": "PK"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":val": (types.AttributeValue)(&attrValue)},
	}
	output, err := svc.PutItem(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.Attributes)
}

// it can do part update
func UpdateItem() {
	var attrValue types.AttributeValueMemberS = types.AttributeValueMemberS{Value: "some value"}
	var time types.AttributeValueMemberS = types.AttributeValueMemberS{Value: time.Now().Format("2006/01/02 15:04:05")}
	input := dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &attrValue,
		},
		TableName:           aws.String("some name"),
		UpdateExpression:    aws.String("set Time = :time"),
		ConditionExpression: aws.String("PK = :val"),

		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val":  &attrValue,
			":time": &time,
		},
		ReturnValues: types.ReturnValueAllNew,
	}
	output, err := svc.UpdateItem(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for key, value := range output.Attributes {
		fmt.Println(key, value)
	}
}

func DeleteItem() {
	var attrValue types.AttributeValueMemberS = types.AttributeValueMemberS{Value: "some value"}
	input := dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &attrValue,
		},
		TableName:    aws.String("some name"),
		ReturnValues: types.ReturnValueAllOld,
	}
	output, err := svc.DeleteItem(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for key, value := range output.Attributes {
		fmt.Println(key, value)
	}
}
