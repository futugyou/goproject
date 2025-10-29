package domain

type StorageMedia interface {
	privateStorageMedia()
	String() string
}

type storageMedia string

func (c storageMedia) privateStorageMedia() {}

func (c storageMedia) String() string {
	return string(c)
}

const (
	StorageMediaLocal      storageMedia = "Local"
	StorageMediaAws        storageMedia = "AWS"
	StorageMediaHCP        storageMedia = "HCP"
	StorageMediaAzureVault storageMedia = "AzureVault"
)

func GetStorageMedia(rType string) StorageMedia {
	switch rType {
	case "Local":
		return StorageMediaLocal
	case "AWS":
		return StorageMediaAws
	case "HCP":
		return StorageMediaHCP
	case "AzureVault":
		return StorageMediaAzureVault
	default:
		return StorageMediaLocal
	}
}
