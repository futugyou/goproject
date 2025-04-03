package configuration

type McpServerConfig struct {
	Id               string
	Name             string
	TransportType    string
	Location         string
	TransportOptions map[string]string
}
