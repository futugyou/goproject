package pipeline

import (
	"errors"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/core/filesystem"
)

type SimpleQueuesConfig struct {
	StorageType                 filesystem.FileSystemTypes
	Directory                   string
	PollDelayMsecs              int
	DispatchFrequencyMsecs      int
	FetchBatchSize              int
	FetchLockSeconds            int
	MaxRetriesBeforePoisonQueue int
	PoisonQueueSuffix           string
}

func NewVolatileConfig() *SimpleQueuesConfig {
	return &SimpleQueuesConfig{
		StorageType:                 filesystem.Volatile,
		Directory:                   "tmp-memory-queues",
		PollDelayMsecs:              100,
		DispatchFrequencyMsecs:      100,
		FetchBatchSize:              3,
		FetchLockSeconds:            300,
		MaxRetriesBeforePoisonQueue: 1,
		PoisonQueueSuffix:           "-poison",
	}
}

func NewPersistentConfig() *SimpleQueuesConfig {
	cfg := NewVolatileConfig()
	cfg.StorageType = filesystem.Disk
	return cfg
}

func (cfg *SimpleQueuesConfig) Validate() error {
	if strings.TrimSpace(cfg.Directory) == "" || strings.Contains(cfg.Directory, " ") {
		return errors.New("SimpleQueue: Directory cannot be empty or have spaces")
	}
	if strings.TrimSpace(cfg.PoisonQueueSuffix) == "" || strings.Contains(cfg.PoisonQueueSuffix, " ") {
		return errors.New("SimpleQueue: PoisonQueueSuffix cannot be empty or have spaces")
	}
	if cfg.PollDelayMsecs < 1 {
		return errors.New("SimpleQueue: PollDelayMsecs cannot be less than 1")
	}
	if cfg.DispatchFrequencyMsecs < 1 {
		return errors.New("SimpleQueue: DispatchFrequencyMsecs cannot be less than 1")
	}
	if cfg.FetchBatchSize < 1 {
		return errors.New("SimpleQueue: FetchBatchSize cannot be less than 1")
	}
	if cfg.FetchLockSeconds < 1 {
		return errors.New("SimpleQueue: FetchLockSeconds cannot be less than 1")
	}
	if cfg.MaxRetriesBeforePoisonQueue < 0 {
		return errors.New("SimpleQueue: MaxRetriesBeforePoisonQueue cannot be less than 0")
	}
	return nil
}
