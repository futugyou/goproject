package client

import "time"

type StdioClientTransportOptions struct {
	Command              string
	Arguments            *string
	WorkingDirectory     *string
	EnvironmentVariables map[string]string
	ShutdownTimeout      time.Duration
	Name                 *string
}

func NewStdioClientTransportOptions(command string) *StdioClientTransportOptions {
	return &StdioClientTransportOptions{
		Command:              command,
		EnvironmentVariables: map[string]string{},
		ShutdownTimeout:      DefaultShutdownTimeout,
	}
}

var DefaultShutdownTimeout time.Duration = 5 * time.Second
