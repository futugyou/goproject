package abstractions

import (
	"fmt"

	"github.com/futugyou/yomawari/core"
)

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

func (c *KernelPluginCollection) GetFunction(pluginName string, functionName string) (KernelFunction, error) {
	result, ok := c.Dictionary.Get(pluginName)
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", pluginName)
	}
	return result.GetFunction(functionName)
}
