package domain

type VaultState interface {
	privateVaultState()
	String() string
}

type vaultState string

func (c vaultState) privateVaultState() {}

func (c vaultState) String() string {
	return string(c)
}

const (
	VaultStateDefault  vaultState = "default"
	VaultStateChanging vaultState = "changing"
)

func GetVaultState(rType string) VaultState {
	switch rType {
	case "default":
		return VaultStateDefault
	case "changing":
		return VaultStateChanging
	default:
		return VaultStateDefault
	}
}
