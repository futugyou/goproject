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
	VaultTypeSystem          vaultType = "system"
	VaultTypeCommon          vaultType = "common"
	VaultTypeProject         vaultType = "project"
	VaultTypeResource        vaultType = "resource"
	VaultTypePlatform        vaultType = "platform"
	VaultTypePlatformProject vaultType = "platform_project"
	VaultTypePlatformWebhook vaultType = "platform_webhook"
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
	case "platform_project":
		return VaultTypePlatformProject
	case "platform_webhook":
		return VaultTypePlatformWebhook
	default:
		return VaultTypeCommon
	}
}
