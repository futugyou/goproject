package abstractions

type KernelMemoryBuilderBuildOptions struct {
	AllowMixingVolatileAndPersistentData bool
}

func NewKernelMemoryBuilderBuildOptions() *KernelMemoryBuilderBuildOptions {
	return &KernelMemoryBuilderBuildOptions{
		AllowMixingVolatileAndPersistentData: false,
	}
}

func WithVolatileAndPersistentData() *KernelMemoryBuilderBuildOptions {
	return &KernelMemoryBuilderBuildOptions{
		AllowMixingVolatileAndPersistentData: true,
	}
}
