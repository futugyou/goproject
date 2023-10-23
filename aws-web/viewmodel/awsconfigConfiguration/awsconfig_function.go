package awsconfigConfiguration

type FunctionConfiguration struct {
	FunctionName         string               `json:"functionName"`
	FunctionArn          string               `json:"functionArn"`
	Runtime              string               `json:"runtime"`
	Role                 string               `json:"role"`
	Handler              string               `json:"handler"`
	CodeSize             int64                `json:"codeSize"`
	Description          string               `json:"description"`
	Timeout              int64                `json:"timeout"`
	MemorySize           int64                `json:"memorySize"`
	LastModified         string               `json:"lastModified"`
	CodeSha256           string               `json:"codeSha256"`
	Version              string               `json:"version"`
	VpcConfig            VpcConfig            `json:"vpcConfig"`
	TracingConfig        TracingConfig        `json:"tracingConfig"`
	RevisionID           string               `json:"revisionId"`
	Layers               []interface{}        `json:"layers"`
	State                string               `json:"state"`
	LastUpdateStatus     string               `json:"lastUpdateStatus"`
	FileSystemConfigs    []FileSystemConfig   `json:"fileSystemConfigs"`
	PackageType          string               `json:"packageType"`
	Architectures        []string             `json:"architectures"`
	EphemeralStorage     EphemeralStorage     `json:"ephemeralStorage"`
	SnapStart            SnapStart            `json:"snapStart"`
	RuntimeVersionConfig RuntimeVersionConfig `json:"runtimeVersionConfig"`
}

type EphemeralStorage struct {
	Size int64 `json:"size"`
}

type FileSystemConfig struct {
	Arn            string `json:"arn"`
	LocalMountPath string `json:"localMountPath"`
}

type RuntimeVersionConfig struct {
	RuntimeVersionArn string `json:"runtimeVersionArn"`
}

type SnapStart struct {
	ApplyOn            string `json:"applyOn"`
	OptimizationStatus string `json:"optimizationStatus"`
}

type TracingConfig struct {
	Mode string `json:"mode"`
}

type VpcConfig struct {
	SubnetIDS        []string `json:"subnetIds"`
	SecurityGroupIDS []string `json:"securityGroupIds"`
}
