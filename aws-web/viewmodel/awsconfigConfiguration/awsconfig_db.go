package awsconfigConfiguration

type DBInstanceConfiguration struct {
	DBInstanceIdentifier                   string                  `json:"dBInstanceIdentifier"`
	DBInstanceClass                        string                  `json:"dBInstanceClass"`
	Engine                                 string                  `json:"engine"`
	DBInstanceStatus                       string                  `json:"dBInstanceStatus"`
	MasterUsername                         string                  `json:"masterUsername"`
	Endpoint                               Endpoint                `json:"endpoint"`
	AllocatedStorage                       int64                   `json:"allocatedStorage"`
	InstanceCreateTime                     string                  `json:"instanceCreateTime"`
	PreferredBackupWindow                  string                  `json:"preferredBackupWindow"`
	BackupRetentionPeriod                  int64                   `json:"backupRetentionPeriod"`
	DBSecurityGroups                       []interface{}           `json:"dBSecurityGroups"`
	VpcSecurityGroups                      []VpcSecurityGroup      `json:"vpcSecurityGroups"`
	AvailabilityZone                       string                  `json:"availabilityZone"`
	DBSubnetGroup                          DBSubnetGroup           `json:"dBSubnetGroup"`
	PreferredMaintenanceWindow             string                  `json:"preferredMaintenanceWindow"`
	PendingModifiedValues                  PendingModifiedValues   `json:"pendingModifiedValues"`
	LatestRestorableTime                   string                  `json:"latestRestorableTime"`
	MultiAZ                                bool                    `json:"multiAZ"`
	EngineVersion                          string                  `json:"engineVersion"`
	AutoMinorVersionUpgrade                bool                    `json:"autoMinorVersionUpgrade"`
	ReadReplicaDBInstanceIdentifiers       []interface{}           `json:"readReplicaDBInstanceIdentifiers"`
	ReadReplicaDBClusterIdentifiers        []interface{}           `json:"readReplicaDBClusterIdentifiers"`
	LicenseModel                           string                  `json:"licenseModel"`
	OptionGroupMemberships                 []OptionGroupMembership `json:"optionGroupMemberships"`
	PubliclyAccessible                     bool                    `json:"publiclyAccessible"`
	StatusInfos                            []interface{}           `json:"statusInfos"`
	StorageType                            string                  `json:"storageType"`
	DBInstancePort                         int64                   `json:"dbInstancePort"`
	StorageEncrypted                       bool                    `json:"storageEncrypted"`
	KmsKeyID                               string                  `json:"kmsKeyId"`
	DbiResourceID                          string                  `json:"dbiResourceId"`
	CACertificateIdentifier                string                  `json:"cACertificateIdentifier"`
	DomainMemberships                      []interface{}           `json:"domainMemberships"`
	CopyTagsToSnapshot                     bool                    `json:"copyTagsToSnapshot"`
	MonitoringInterval                     int64                   `json:"monitoringInterval"`
	EnhancedMonitoringResourceArn          string                  `json:"enhancedMonitoringResourceArn"`
	MonitoringRoleArn                      string                  `json:"monitoringRoleArn"`
	DBInstanceArn                          string                  `json:"dBInstanceArn"`
	IAMDatabaseAuthenticationEnabled       bool                    `json:"iAMDatabaseAuthenticationEnabled"`
	PerformanceInsightsEnabled             bool                    `json:"performanceInsightsEnabled"`
	EnabledCloudwatchLogsExports           []interface{}           `json:"enabledCloudwatchLogsExports"`
	ProcessorFeatures                      []interface{}           `json:"processorFeatures"`
	DeletionProtection                     bool                    `json:"deletionProtection"`
	AssociatedRoles                        []interface{}           `json:"associatedRoles"`
	MaxAllocatedStorage                    int64                   `json:"maxAllocatedStorage"`
	TagList                                []interface{}           `json:"tagList"`
	DBInstanceAutomatedBackupsReplications []interface{}           `json:"dBInstanceAutomatedBackupsReplications"`
	CustomerOwnedIPEnabled                 bool                    `json:"customerOwnedIpEnabled"`
}

type DBSubnetGroup struct {
	DBSubnetGroupName        string   `json:"dBSubnetGroupName"`
	DBSubnetGroupDescription string   `json:"dBSubnetGroupDescription"`
	VpcID                    string   `json:"vpcId"`
	SubnetGroupStatus        string   `json:"subnetGroupStatus"`
	Subnets                  []Subnet `json:"subnets"`
}

type SubnetAvailabilityZone struct {
	Name string `json:"name"`
}

type Endpoint struct {
	Address      string `json:"address"`
	Port         int64  `json:"port"`
	HostedZoneID string `json:"hostedZoneId"`
}

type OptionGroupMembership struct {
	OptionGroupName string `json:"optionGroupName"`
	Status          string `json:"status"`
}

type PendingModifiedValues struct {
	ProcessorFeatures []interface{} `json:"processorFeatures"`
}

type VpcSecurityGroup struct {
	VpcSecurityGroupID string `json:"vpcSecurityGroupId"`
	Status             string `json:"status"`
}

type SecurityGroupConfiguration struct {
	Description         string         `json:"description"`
	GroupName           string         `json:"groupName"`
	IPPermissions       []IPPermission `json:"ipPermissions"`
	OwnerID             string         `json:"ownerId"`
	GroupID             string         `json:"groupId"`
	IPPermissionsEgress []IPPermission `json:"ipPermissionsEgress"`
	Tags                []interface{}  `json:"tags"`
	VpcID               string         `json:"vpcId"`
}

type IPPermission struct {
	IPProtocol       string            `json:"ipProtocol"`
	Ipv6Ranges       []interface{}     `json:"ipv6Ranges"`
	PrefixListIDS    []interface{}     `json:"prefixListIds"`
	UserIDGroupPairs []UserIDGroupPair `json:"userIdGroupPairs"`
	Ipv4Ranges       []Ipv4Range       `json:"ipv4Ranges"`
	IPRanges         []string          `json:"ipRanges"`
}

type Ipv4Range struct {
	CIDRIP string `json:"cidrIp"`
}

type UserIDGroupPair struct {
	GroupID string `json:"groupId"`
	UserID  string `json:"userId"`
}
