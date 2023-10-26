package awsconfigConfiguration

type ECSTaskDefinitionConfiguration struct {
	Status                  string                `json:"Status"`
	TaskRoleArn             string                `json:"TaskRoleArn"`
	InferenceAccelerators   []interface{}         `json:"InferenceAccelerators"`
	Memory                  string                `json:"Memory"`
	PlacementConstraints    []interface{}         `json:"PlacementConstraints"`
	CPU                     string                `json:"Cpu"`
	RequiresCompatibilities []string              `json:"RequiresCompatibilities"`
	NetworkMode             string                `json:"NetworkMode"`
	ExecutionRoleArn        string                `json:"ExecutionRoleArn"`
	Volumes                 []Volume              `json:"Volumes"`
	ContainerDefinitions    []ContainerDefinition `json:"ContainerDefinitions"`
	Family                  string                `json:"Family"`
	Tags                    []Tag                 `json:"Tags"`
	TaskDefinitionArn       string                `json:"TaskDefinitionArn"`
}

type ContainerDefinition struct {
	Secrets               []Secret         `json:"Secrets"`
	ExtraHosts            []interface{}    `json:"ExtraHosts"`
	VolumesFrom           []interface{}    `json:"VolumesFrom"`
	CPU                   float64          `json:"Cpu"`
	EntryPoint            []interface{}    `json:"EntryPoint"`
	DNSServers            []interface{}    `json:"DnsServers"`
	Image                 string           `json:"Image"`
	Essential             bool             `json:"Essential"`
	LogConfiguration      LogConfiguration `json:"LogConfiguration"`
	ResourceRequirements  []interface{}    `json:"ResourceRequirements"`
	EnvironmentFiles      []interface{}    `json:"EnvironmentFiles"`
	Name                  string           `json:"Name"`
	MountPoints           []MountPoint     `json:"MountPoints"`
	DependsOn             []DependsOn      `json:"DependsOn"`
	PortMappings          []PortMapping    `json:"PortMappings"`
	DockerLabels          DockerLabels     `json:"DockerLabels"`
	DockerSecurityOptions []interface{}    `json:"DockerSecurityOptions"`
	SystemControls        []interface{}    `json:"SystemControls"`
	Command               []interface{}    `json:"Command"`
	DNSSearchDomains      []interface{}    `json:"DnsSearchDomains"`
	Links                 []interface{}    `json:"Links"`
	Environment           []Environment    `json:"Environment"`
	Ulimits               []interface{}    `json:"Ulimits"`
	CredentialSpecs       []interface{}    `json:"CredentialSpecs"`
}

type DependsOn struct {
	Condition     string `json:"Condition"`
	ContainerName string `json:"ContainerName"`
}

type DockerLabels struct {
}

type Environment struct {
	Value string `json:"Value"`
	Name  string `json:"Name"`
}

type LogConfiguration struct {
	SecretOptions []interface{} `json:"SecretOptions"`
	Options       Options       `json:"Options"`
	LogDriver     string        `json:"LogDriver"`
}

type Options struct {
	AwslogsGroup        string `json:"awslogs-group"`
	AwslogsRegion       string `json:"awslogs-region"`
	AwslogsStreamPrefix string `json:"awslogs-stream-prefix"`
}

type MountPoint struct {
	SourceVolume  string `json:"SourceVolume"`
	ContainerPath string `json:"ContainerPath"`
}

type PortMapping struct {
	HostPort      float64 `json:"HostPort"`
	ContainerPort float64 `json:"ContainerPort"`
	Protocol      string  `json:"Protocol"`
}

type Secret struct {
	ValueFrom string `json:"ValueFrom"`
	Name      string `json:"Name"`
}

type Volume struct {
	EFSVolumeConfiguration EFSVolumeConfiguration `json:"EfsVolumeConfiguration"`
	Name                   string                 `json:"Name"`
}

type EFSVolumeConfiguration struct {
	TransitEncryption   string              `json:"TransitEncryption"`
	AuthorizationConfig AuthorizationConfig `json:"AuthorizationConfig"`
	FileSystemID        string              `json:"FileSystemId"`
	RootDirectory       string              `json:"RootDirectory"`
}

type AuthorizationConfig struct {
	AccessPointID string `json:"AccessPointId"`
}
