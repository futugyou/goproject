package services

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"

	"log"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
)

// path for aws config snapshot data (download from s3)
func (a *AwsConfigService) SyncFileResources(path string) {
	// 1. read data from file
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	var rawDatas []model.AwsConfigRawData

	json.Unmarshal(byteValue, &rawDatas)

	if len(rawDatas) == 0 {
		return
	}

	a.commonDataExec(rawDatas)
}

func (a *AwsConfigService) SyncResourcesByConfig() {
	// 1. init account
	if !initAwsEnv() {
		return
	}

	// 2. get data from aws config
	rawdataStringList := getAwsConfigData()
	if len(rawdataStringList) == 0 {
		return
	}

	// 3. convert raw data to file data
	rawData := convertToRawData(rawdataStringList)
	if len(rawData) == 0 {
		return
	}
	a.commonDataExec(rawData)
}

func convertToRawData(rawdata []string) []model.AwsConfigRawData {
	rawDatas := make([]model.AwsConfigRawData, 0)

	rawdatastring := "[" + strings.Join(rawdata, ",") + "]"
	err := json.Unmarshal([]byte(rawdatastring), &rawDatas)
	if err != nil {
		log.Println(err)
		return rawDatas
	}
	return rawDatas
}

func initAwsEnv() bool {
	accountService := NewAccountService()
	accountid := os.Getenv("accountid")
	if len(accountid) == 0 {
		log.Println("can not find accountid from env.")
		return false
	}

	account := accountService.GetAccountByID(accountid)
	if account == nil {
		log.Printf("can not find accountid:%s from db.", accountid)
		return false
	}

	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	return true
}

func getAwsConfigData() []string {
	svc := configservice.NewFromConfig(awsenv.Cfg)
	var nextToken *string = nil
	results := make([]string, 0)
	for {
		input := &configservice.SelectResourceConfigInput{
			Expression: aws.String(`
	SELECT
		version,
		accountId,
		configurationItemCaptureTime,
		configurationItemStatus,
		configurationStateId,
		arn,
		resourceType,
		resourceId,
		resourceName,
		awsRegion,
		availabilityZone,
		tags,
		relatedEvents,
		relationships,
		configuration,
		supplementaryConfiguration,
		resourceTransitionStatus,
		resourceCreationTime
	WHERE
		resourceType <> 'AWS::Backup::RecoveryPoint'
		and resourceType <> 'AWS::CodeDeploy::DeploymentConfig'
		and resourceType <> 'AWS::RDS::DBSnapshot'
  `),
			Limit:     100,
			NextToken: nextToken,
		}
		log.Println(1)
		output, err := svc.SelectResourceConfig(context.Background(), input)
		if err != nil {
			log.Println(err)
			return []string{}
		}
		log.Println(2)
		results = append(results, output.Results...)
		nextToken = output.NextToken

		if output.NextToken == nil {
			break
		}
	}

	return results
}

func (a *AwsConfigService) commonDataExec(rawDatas []model.AwsConfigRawData) {
	// 2. filter data
	rawDatas = FilterResource(rawDatas)

	// 3. get all vpc info
	vpcinfos := GetAllVpcInfos(rawDatas)

	resources := make([]entity.AwsConfigEntity, 0)
	ships := make([]entity.AwsConfigRelationshipEntity, 0)

	// 4. create AwsConfigEntity list
	for _, data := range rawDatas {
		resource := CreateAwsConfigEntity(data, vpcinfos)
		resources = append(resources, resource)
	}

	// 4.1 add individual resource
	resources = AddIndividualResource(resources, vpcinfos)

	// 5. create AwsConfigRelationshipEntity list
	for _, data := range rawDatas {
		ship := CreateAwsConfigRelationshipEntity(data, resources)
		ships = append(ships, ship...)
	}

	// 5.1 individual relation ship
	individualShips := AddIndividualRelationShip(resources)
	ships = append(ships, individualShips...)

	// 5.2 remove duplicate
	ships = RemoveDuplicateRelationShip(ships)

	log.Println("resources count: ", len(resources))
	log.Println("relationships count: ", len(ships))
	if len(resources) == 0 || len(ships) == 0 {
		return
	}

	// 6. delete all
	err := a.repository.DeleteAll(context.Background())
	if err != nil {
		log.Println("delete awsconfig error: ", err.Error())
		return
	}
	err = a.relRepository.DeleteAll(context.Background())
	if err != nil {
		log.Println("delete awsconfigrelationship error: ", err.Error())
		return
	}

	// 7. Insert all
	err = a.repository.InsertMany(context.Background(), resources)
	if err != nil {
		log.Println("Insert awsconfig error: ", err.Error())
		return
	}
	err = a.relRepository.InsertMany(context.Background(), ships)
	if err != nil {
		log.Println("Insert awsconfigrelationship error: ", err.Error())
	}
}
