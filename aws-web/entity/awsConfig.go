package entity

import "time"

type AwsConfigEntity struct {
	ID                           string    `bson:"_id"`
	Label                        string    `bson:"label"`
	AccountID                    string    `bson:"accountId"`
	Arn                          string    `bson:"arn"`
	AvailabilityZone             string    `bson:"availabilityZone"`
	AwsRegion                    string    `bson:"awsRegion"`
	Configuration                string    `bson:"configuration"`
	ConfigurationItemCaptureTime time.Time `bson:"configurationItemCaptureTime"`
	ConfigurationItemStatus      string    `bson:"configurationItemStatus"`
	ConfigurationStateID         int64     `bson:"configurationStateId"`
	ResourceCreationTime         time.Time `bson:"resourceCreationTime"`
	ResourceID                   string    `bson:"resourceId"`
	ResourceName                 string    `bson:"resourceName"`
	ResourceType                 string    `bson:"resourceType"`
	Tags                         string    `bson:"tags"`
	Version                      string    `bson:"version"`
	VpcID                        string    `bson:"vpcId"`
	SubnetID                     string    `bson:"subnetId"`
	SubnetIds                    []string  `bson:"subnetIds"`
	Title                        string    `bson:"title"`
}

func (AwsConfigEntity) GetType() string {
	return "awsConfigs"
}

type AwsConfigRelationshipEntity struct {
	ID          string `bson:"_id"`
	Label       string `bson:"label"`
	SourceID    string `bson:"sourceId"`
	SourceLabel string `bson:"sourceLabel"`
	TargetID    string `bson:"targetId"`
	TargetLabel string `bson:"targetLabel"`
}

func (AwsConfigRelationshipEntity) GetType() string {
	return "awsConfigRelationships"
}
