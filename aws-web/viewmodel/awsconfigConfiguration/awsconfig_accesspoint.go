package awsconfigConfiguration

type AccessPointConfiguration struct {
	AccessPointID   string        `json:"AccessPointId"`
	Arn             string        `json:"Arn"`
	ClientToken     string        `json:"ClientToken"`
	AccessPointTags []Tag         `json:"AccessPointTags"`
	FileSystemID    string        `json:"FileSystemId"`
	POSIXUser       POSIXUser     `json:"PosixUser"`
	RootDirectory   RootDirectory `json:"RootDirectory"`
}

type POSIXUser struct {
	Uid           string   `json:"Uid"`
	Gid           string   `json:"Gid"`
	SecondaryGids []string `json:"SecondaryGids"`
}

type RootDirectory struct {
	Path         string       `json:"Path"`
	CreationInfo CreationInfo `json:"CreationInfo"`
}

type CreationInfo struct {
	OwnerUid    string `json:"OwnerUid"`
	OwnerGid    string `json:"OwnerGid"`
	Permissions string `json:"Permissions"`
}
