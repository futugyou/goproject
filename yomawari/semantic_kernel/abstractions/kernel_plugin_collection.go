package abstractions

import "github.com/futugyou/yomawari/core"

type KernelPluginCollection struct {
	*core.Dictionary[string, KernelPlugin]
}

func NewKernelPluginCollection(plugins []KernelPlugin) *KernelPluginCollection {
	d := core.NewDictionary[string, KernelPlugin]()
	for _, plugin := range plugins {
		d.Set(plugin.Name(), plugin)
	}
	return &KernelPluginCollection{
		Dictionary: d,
	}
}

func (c *KernelPluginCollection) Get(name string) KernelPlugin {
	result, _ := c.Dictionary.Get(name)
	return result
}

func (c *KernelPluginCollection) Add(plugin KernelPlugin) {
	c.Set(plugin.Name(), plugin)
}

func (c *KernelPluginCollection) AddRange(plugins []KernelPlugin) {
	for _, plugin := range plugins {
		c.Add(plugin)
	}
}

func (c *KernelPluginCollection) Remove(plugin KernelPlugin) {
	c.Dictionary.Remove(plugin.Name())
}
