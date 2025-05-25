package abstractions

type Kernel struct {
	data    map[string]any
	plugins KernelPluginCollection
}
