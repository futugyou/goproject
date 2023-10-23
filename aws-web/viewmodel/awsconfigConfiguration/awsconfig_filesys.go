package awsconfigConfiguration

type FileSystemConfiguration struct {
	FileSystemID      string            `json:"FileSystemId"`
	Arn               string            `json:"Arn"`
	Encrypted         bool              `json:"Encrypted"`
	FileSystemTags    []FileSystemTag   `json:"FileSystemTags"`
	PerformanceMode   string            `json:"PerformanceMode"`
	ThroughputMode    string            `json:"ThroughputMode"`
	LifecyclePolicies []LifecyclePolicy `json:"LifecyclePolicies"`
	BackupPolicy      BackupPolicy      `json:"BackupPolicy"`
	FileSystemPolicy  FileSystemPolicy  `json:"FileSystemPolicy"`
	KmsKeyID          string            `json:"KmsKeyId"`
}

type BackupPolicy struct {
	Status string `json:"Status"`
}

type FileSystemPolicy struct {
}

type FileSystemTag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type LifecyclePolicy struct {
	TransitionToIA                  *string `json:"TransitionToIA,omitempty"`
	TransitionToPrimaryStorageClass *string `json:"TransitionToPrimaryStorageClass,omitempty"`
}
