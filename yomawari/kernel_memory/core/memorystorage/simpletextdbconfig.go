package memorystorage

import "github.com/futugyou/yomawari/kernel_memory/core/filesystem"

type SimpleTextDbConfig struct {
	StorageType filesystem.FileSystemTypes
	Directory   string
}

func NewSimpleTextDbConfig() *SimpleTextDbConfig {
	return &SimpleTextDbConfig{
		StorageType: filesystem.Volatile,
		Directory:   "tmp-memory-text",
	}
}
func VolatileSimpleTextDbConfig() *SimpleTextDbConfig {
	return &SimpleTextDbConfig{
		StorageType: filesystem.Volatile,
		Directory:   "tmp-memory-text",
	}
}
func DiskSimpleTextDbConfig() *SimpleTextDbConfig {
	return &SimpleTextDbConfig{
		StorageType: filesystem.Disk,
		Directory:   "tmp-memory-text",
	}
}
