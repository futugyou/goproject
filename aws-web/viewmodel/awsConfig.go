package viewmodel

import (
	"time"
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
