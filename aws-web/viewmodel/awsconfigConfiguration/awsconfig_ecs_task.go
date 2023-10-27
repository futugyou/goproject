package awsconfigConfiguration

type ECSTaskConfiguration struct {
	Attachments          []EcsAttachment `json:"Attachments"`
	AvailabilityZone     string          `json:"AvailabilityZone"`
	CapacityProviderName string          `json:"CapacityProviderName"`
	ClusterArn           string          `json:"ClusterArn"`
	ContainerInstanceArn string          `json:"ContainerInstanceArn"`
	Containers           []Container     `json:"Containers"`
	Cpu                  string          `json:"Cpu"`
	Group                string          `json:"Group"`
	HealthStatus         string          `json:"HealthStatus"`
	LastStatus           string          `json:"LastStatus"`
	LaunchType           string          `json:"LaunchType"`
	Memory               string          `json:"Memory"`
	Overrides            TaskOverride    `json:"Overrides"`
	PlatformFamily       string          `json:"PlatformFamily"`
	PlatformVersion      string          `json:"PlatformVersion"`
	Tags                 []Tag           `json:"Tags"`
	TaskArn              string          `json:"TaskArn"`
	TaskDefinitionArn    string          `json:"TaskDefinitionArn"`
	Version              int64           `json:"Version"`
}
type EcsAttachment struct {
	Details []KeyValuePair `json:"Details"`
	Id      string         `json:"Id"`
	Status  string         `json:"Status"`
	Type    string         `json:"Type"`
}

type TaskOverride struct {
	ContainerOverrides            []ContainerOverride            `json:"ContainerOverrides"`
	Cpu                           string                         `json:"Cpu"`
	EphemeralStorage              EphemeralStorage               `json:"EphemeralStorage"`
	ExecutionRoleArn              string                         `json:"ExecutionRoleArn"`
	InferenceAcceleratorOverrides []InferenceAcceleratorOverride `json:"InferenceAcceleratorOverrides"`
	Memory                        string                         `json:"Memory"`
	TaskRoleArn                   string                         `json:"TaskRoleArn"`
}

type ContainerOverride struct {
	Command           []string       `json:"Command"`
	Cpu               int32          `json:"Cpu"`
	Environment       []KeyValuePair `json:"Environment"`
	Memory            int32          `json:"Memory"`
	MemoryReservation int32          `json:"MemoryReservation"`
	Name              string         `json:"Name"`
}

type InferenceAcceleratorOverride struct {
	DeviceName string `json:"DeviceName"`
	DeviceType string `json:"DeviceType"`
}

type KeyValuePair struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type Container struct {
	ContainerArn      string                `json:"ContainerArn"`
	Cpu               string                `json:"Cpu"`
	ExitCode          int32                 `json:"ExitCode"`
	GpuIds            []string              `json:"GpuIds"`
	Image             string                `json:"Image"`
	ImageDigest       string                `json:"ImageDigest"`
	LastStatus        string                `json:"LastStatus"`
	Memory            string                `json:"Memory"`
	MemoryReservation string                `json:"MemoryReservation"`
	Name              string                `json:"Name"`
	NetworkInterfaces []EcsNetworkInterface `json:"NetworkInterfaces"`
	Reason            string                `json:"Reason"`
	RuntimeId         string                `json:"RuntimeId"`
	TaskArn           string                `json:"TaskArn"`
}

type EcsNetworkInterface struct {
	AttachmentId       string `json:"AttachmentId"`
	Ipv6Address        string `json:"Ipv6Address"`
	PrivateIpv4Address string `json:"PrivateIpv4Address"`
}
