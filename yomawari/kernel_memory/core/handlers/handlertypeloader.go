package handlers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/core/configuration"
)

func TryGetHandlerType(config configuration.HandlerConfig) (interface{}, error) {
	if config.Class == "" || config.Assembly == "" {
		return nil, errors.New("handler type loader: class or assembly is empty")
	}

	assemblyPaths := []string{
		config.Assembly,
		filepath.Join(".", config.Assembly),
		filepath.Join(filepath.Dir(os.Args[0]), config.Assembly),
	}

	var path string
	for _, p := range assemblyPaths {
		if _, err := os.Stat(p); err == nil {
			path = p
			break
		}
	}

	if path == "" {
		return nil, fmt.Errorf("handler type loader: handler assembly not found: %s", config.Assembly)
	}

	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %v", err)
	}

	sym, err := p.Lookup(config.Class)
	if err != nil {
		return nil, fmt.Errorf("handler type loader: invalid handler definition: `%s` not found", config.Class)
	}

	handler, ok := sym.(pipeline.IPipelineStepHandler)
	if !ok {
		return nil, fmt.Errorf("handler type loader: `%s` does not implement IPipelineStepHandler", config.Class)
	}

	return handler, nil
}
