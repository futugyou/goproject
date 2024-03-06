package dynamo2struct

import (
	"fmt"
)

type DynamoDBConfig struct {
	EntityFolder     string
	RepoFolder       string
	CoreFoler        string
	DynamoRepoFolder string
	PkgName          string
	AccessKey        string
	AccessSecret     string
	Region           string
}

func (m *DynamoDBConfig) Check() error {
	if len(m.AccessKey) == 0 {
		return fmt.Errorf("dynamodb AccessKey can not be nil")
	}
	if len(m.AccessSecret) == 0 {
		return fmt.Errorf("dynamodb AccessSecret can not be nil")
	}
	if len(m.Region) == 0 {
		return fmt.Errorf("dynamodb Region can not be nil")
	}
	return nil
}
