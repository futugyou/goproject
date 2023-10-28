package awsconfigConfiguration

type UserConfiguration struct {
	Path                    string        `json:"path"`
	AttachedManagedPolicies []Policy      `json:"attachedManagedPolicies"`
	GroupList               []string      `json:"groupList"`
	UserPolicyList          []interface{} `json:"userPolicyList"`
	UserName                string        `json:"userName"`
	Arn                     string        `json:"arn"`
	UserID                  string        `json:"userId"`
	CreateDate              string        `json:"createDate"`
	Tags                    []Tag         `json:"tags"`
}
