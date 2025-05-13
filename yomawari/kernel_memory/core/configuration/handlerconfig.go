package configuration

type HandlerConfig struct {
	Assembly string
	Class    string
}

func NewHandlerConfig(assembly string, class string) *HandlerConfig {
	return &HandlerConfig{
		Assembly: assembly,
		Class:    class,
	}
}

const DocsBaseUrl string = "https://microsoft.github.io/kernel_memory"
