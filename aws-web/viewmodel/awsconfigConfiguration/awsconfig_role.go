package awsconfigConfiguration

type RoleConfiguration struct {
	Path                     string        `json:"path"`
	AssumeRolePolicyDocument string        `json:"assumeRolePolicyDocument"`
	InstanceProfileList      []interface{} `json:"instanceProfileList"`
	RoleID                   string        `json:"roleId"`
	AttachedManagedPolicies  []Policy      `json:"attachedManagedPolicies"`
	RoleName                 string        `json:"roleName"`
	Arn                      string        `json:"arn"`
	CreateDate               string        `json:"createDate"`
	RolePolicyList           []Policy      `json:"rolePolicyList"`
	Tags                     []Tag         `json:"tags"`
}

type Policy struct {
	PolicyArn      string `json:"policyArn"`
	PolicyName     string `json:"policyName"`
	PolicyDocument string `json:"policyDocument"`
}
