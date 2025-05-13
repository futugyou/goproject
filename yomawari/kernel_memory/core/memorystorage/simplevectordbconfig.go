package memorystorage

import "github.com/futugyou/yomawari/kernel_memory/core/filesystem"

type SimpleVectorDbConfig struct {
	StorageType filesystem.FileSystemTypes
	Directory   string
}

func NewSimpleVectorDbConfig() *SimpleVectorDbConfig {
	return &SimpleVectorDbConfig{
		StorageType: filesystem.Volatile,
		Directory:   "tmp-memory-vectors",
	}
}
func VolatileSimpleVectorDbConfig() *SimpleVectorDbConfig {
	return &SimpleVectorDbConfig{
		StorageType: filesystem.Volatile,
		Directory:   "tmp-memory-vectors",
	}
}
func DiskSimpleVectorDbConfig() *SimpleVectorDbConfig {
	return &SimpleVectorDbConfig{
		StorageType: filesystem.Disk,
		Directory:   "tmp-memory-vectors",
	}
}
