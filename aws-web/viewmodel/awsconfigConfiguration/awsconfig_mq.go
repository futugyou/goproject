package awsconfigConfiguration

type AmazonMQConfiguration struct {
	SecurityGroups          []string `json:"SecurityGroups"`
	SubnetIDS               []string `json:"SubnetIds"`
	DeploymentMode          string   `json:"DeploymentMode"`
	EngineType              string   `json:"EngineType"`
	Tags                    []Tag    `json:"Tags"`
	ConfigurationRevision   int64    `json:"ConfigurationRevision"`
	StorageType             string   `json:"StorageType"`
	EngineVersion           string   `json:"EngineVersion"`
	HostInstanceType        string   `json:"HostInstanceType"`
	AutoMinorVersionUpgrade bool     `json:"AutoMinorVersionUpgrade"`
}
