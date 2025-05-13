package documentstorage

import (
	"github.com/futugyou/yomawari/kernel_memory/core/filesystem"
)

type SimpleFileStorageConfig struct {
	StorageType filesystem.FileSystemTypes
	Directory   string
}

func Volatile() SimpleFileStorageConfig {
	return SimpleFileStorageConfig{
		StorageType: filesystem.Volatile,
		Directory:   "tmp-memory-files",
	}
}

func Persistent() SimpleFileStorageConfig {
	return SimpleFileStorageConfig{
		StorageType: filesystem.Disk,
		Directory:   "tmp-memory-files",
	}
}
