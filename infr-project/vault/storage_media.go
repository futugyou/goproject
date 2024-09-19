package vault

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
	StorageMediaLocal storageMedia = "Local"
	StorageMediaAws   storageMedia = "AWS"
	StorageMediaHCP   storageMedia = "HCP"
)

func GetStorageMedia(rType string) StorageMedia {
	switch rType {
	case "Local":
		return StorageMediaLocal
	case "AWS":
		return StorageMediaAws
	case "HCP":
		return StorageMediaHCP
	default:
		return StorageMediaLocal
	}
}
