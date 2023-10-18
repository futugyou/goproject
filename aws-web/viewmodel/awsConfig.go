package viewmodel

import (
	"strings"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type ResourceGraph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Edge struct {
	ID     string   `json:"id"`
	Label  string   `json:"label"`
	Source EdgeItem `json:"source"`
	Target EdgeItem `json:"target"`
}

type EdgeItem struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	ResourceType string `json:"resourceType"`
}

type Node struct {
	ID         string     `json:"id"`
	Label      string     `json:"label"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	AccountID                    string    `json:"accountId"`
	Arn                          string    `json:"arn"`
	AvailabilityZone             string    `json:"availabilityZone"`
	AwsRegion                    string    `json:"awsRegion"`
	Configuration                string    `json:"configuration"`
	ConfigurationItemCaptureTime time.Time `json:"configurationItemCaptureTime"`
	ConfigurationItemStatus      string    `json:"configurationItemStatus"`
	ConfigurationStateID         int64     `json:"configurationStateId"`
	ResourceCreationTime         time.Time `json:"resourceCreationTime"`
	ResourceID                   string    `json:"resourceId"`
	ResourceName                 string    `json:"resourceName"`
	ResourceType                 string    `json:"resourceType"`
	Tags                         string    `json:"tags"`
	Version                      string    `json:"version"`
	VpcID                        string    `json:"vpcId"`
	SubnetID                     string    `json:"subnetId"`
	SubnetIDS                    []string  `json:"subnetIds"`
	Title                        string    `json:"title"`
	SecurityGroups               []string  `json:"securityGroups"`
	LoginURL                     string    `json:"loginURL"`
	LoggedInURL                  string    `json:"loggedInURL"`
}

type AwsConfigFileData struct {
	RelatedEvents                []string          `json:"relatedEvents"`
	Relationships                []Relationship    `json:"relationships"`
	Configuration                interface{}       `json:"configuration"`
	Tags                         map[string]string `json:"tags"`
	ConfigurationItemVersion     string            `json:"configurationItemVersion"`
	ConfigurationItemCaptureTime time.Time         `json:"configurationItemCaptureTime"`
	ConfigurationStateID         int64             `json:"configurationStateId"`
	AwsAccountID                 string            `json:"awsAccountId"`
	ConfigurationItemStatus      string            `json:"configurationItemStatus"`
	ResourceType                 string            `json:"resourceType"`
	ResourceID                   string            `json:"resourceId"`
	ResourceName                 string            `json:"resourceName"`
	ARN                          string            `json:"ARN"`
	AwsRegion                    string            `json:"awsRegion"`
	AvailabilityZone             string            `json:"availabilityZone"`
	ConfigurationStateMd5Hash    string            `json:"configurationStateMd5Hash"`
	ResourceCreationTime         time.Time         `json:"resourceCreationTime"`
}

type Relationship struct {
	ResourceID   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
}

func (data AwsConfigFileData) CreateAwsConfigEntity(vpcinfos []VpcInfo) entity.AwsConfigEntity {
	configuration := getDataString(data.Configuration)
	name := getName(data.ResourceID, data.ResourceName, data.Tags)
	vpcid, subnetId, subnetIds, securityGroups := getVpcInfo(data.ResourceType, configuration, vpcinfos)
	loginURL, loggedInURL := createConsoleUrls(data)
	config := entity.AwsConfigEntity{
		ID:                           getId(data.ARN, data.ResourceID),
		Label:                        name,
		AccountID:                    data.AwsAccountID,
		Arn:                          data.ARN,
		AvailabilityZone:             data.AvailabilityZone,
		AwsRegion:                    data.AwsRegion,
		Configuration:                configuration,
		ConfigurationItemCaptureTime: data.ConfigurationItemCaptureTime,
		ConfigurationItemStatus:      data.ConfigurationItemStatus,
		ConfigurationStateID:         data.ConfigurationStateID,
		ResourceCreationTime:         data.ResourceCreationTime,
		ResourceID:                   data.ResourceID,
		ResourceName:                 data.ResourceName,
		ResourceType:                 data.ResourceType,
		Tags:                         getDataString(data.Tags),
		Version:                      data.ConfigurationItemVersion,
		VpcID:                        vpcid,
		SubnetID:                     subnetId,
		SubnetIds:                    subnetIds,
		Title:                        name,
		SecurityGroups:               securityGroups,
		LoginURL:                     loginURL,
		LoggedInURL:                  loggedInURL,
	}
	return config
}

func (data AwsConfigFileData) CreateAwsConfigRelationshipEntity(configs []entity.AwsConfigEntity) []entity.AwsConfigRelationshipEntity {
	lists := make([]entity.AwsConfigRelationshipEntity, 0)

	for _, ship := range data.Relationships {
		var id string
		for i := 0; i < len(configs); i++ {
			if configs[i].ResourceID == ship.ResourceID && configs[i].ResourceName == ship.ResourceName {
				id = getId(configs[i].Arn, configs[i].ResourceID)
				break
			}
		}

		if len(id) == 0 {
			continue
		}

		if strings.HasPrefix(ship.Name, "Is ") {
			relationship := entity.AwsConfigRelationshipEntity{
				ID:                 data.ResourceID + "-" + ship.ResourceID,
				SourceID:           getId(data.ARN, data.ResourceID),
				SourceLabel:        data.ResourceName,
				SourceResourceType: data.ResourceType,
				Label:              ship.Name,
				TargetID:           id,
				TargetLabel:        ship.ResourceName,
				TargetResourceType: ship.ResourceType,
			}
			lists = append(lists, relationship)
		}

		if strings.HasPrefix(ship.Name, "Contains ") {
			relationship := entity.AwsConfigRelationshipEntity{
				ID:                 ship.ResourceID + "-" + data.ResourceID,
				SourceID:           id,
				SourceLabel:        ship.ResourceName,
				SourceResourceType: ship.ResourceType,
				Label:              ship.Name,
				TargetID:           getId(data.ARN, data.ResourceID),
				TargetLabel:        data.ResourceName,
				TargetResourceType: data.ResourceType,
			}
			lists = append(lists, relationship)
		}

	}

	return lists
}
