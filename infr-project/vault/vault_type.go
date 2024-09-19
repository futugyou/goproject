package vault

type VaultType interface {
	privateVaultType()
	String() string
}

type vaultType string

func (c vaultType) privateVaultType() {}

func (c vaultType) String() string {
	return string(c)
}

const (
	VaultTypeSystem   vaultType = "system"
	VaultTypeCommon   vaultType = "common"
	VaultTypeProject  vaultType = "project"
	VaultTypeResource vaultType = "resource"
	VaultTypePlatform vaultType = "platform"
)

func GetVaultType(rType string) VaultType {
	switch rType {
	case "system":
		return VaultTypeSystem
	case "common":
		return VaultTypeCommon
	case "project":
		return VaultTypeProject
	case "resource":
		return VaultTypeResource
	case "platform":
		return VaultTypePlatform
	default:
		return VaultTypeCommon
	}
}
