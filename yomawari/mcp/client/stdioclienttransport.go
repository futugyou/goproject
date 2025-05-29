package client

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/futugyou/yomawari/mcp/configuration"
	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ IClientTransport = (*StdioClientTransport)(nil)

type StdioClientTransport struct {
	options      *StdioClientTransportOptions
	serverConfig *configuration.McpServerConfig
	cmd          *exec.Cmd
	logger       logging.Logger
	name         string
}

// GetName implements IClientTransport.
func (s *StdioClientTransport) GetName() string {
	return s.name
}

func NewStdioClientTransport(serverConfig *configuration.McpServerConfig, options *StdioClientTransportOptions, logger logging.Logger) *StdioClientTransport {
	t := &StdioClientTransport{
		options:      options,
		serverConfig: serverConfig,
		logger:       logger,
	}

	if options != nil && options.Name != nil {
		t.name = *options.Name
	}

	if len(t.name) == 0 {
		t.name = fmt.Sprintf("stdio-%s", strings.ReplaceAll(t.options.Command, " ", "-"))
	}
	return t
}

func convertEnvVars(envVars map[string]string) []string {
	var envs []string
	for k, v := range envVars {
		envs = append(envs, k+"="+v)
	}
	return envs
}

// Connect implements IClientTransport.
func (s *StdioClientTransport) Connect(context.Context) (protocol.ITransport, error) {
	endpointName := fmt.Sprintf("Client (stdio) for (%s: %s)", s.serverConfig.Id, s.serverConfig.Name)
	fmt.Println("Connecting:", endpointName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := strings.Fields(*s.options.Arguments)

	cmd := exec.CommandContext(ctx, s.options.Command, args...)
	cmd.Dir = *s.options.WorkingDirectory
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Env = append(os.Environ(), convertEnvVars(s.options.EnvironmentVariables)...)

	_, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("can not create stdin pipeline: %w", err)
	}
	_, err = cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("can not create stdout pipeline: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("can not create stderr pipeline: %w", err)
	}
	stderr := &bytes.Buffer{}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("can not start process: %w", err)
	}

	go func() {
		_, err := stderr.ReadFrom(stderrPipe)
		if err != nil {
			fmt.Println("read stderr error:", err)
		}
	}()

	fmt.Println("process started, PID:", cmd.Process.Pid)

	s.cmd = cmd
	return NewStdioClientSessionTransport(s.options, cmd, endpointName, s.logger), nil
}

func (t *StdioClientTransport) Dispose(shutdownTimeout time.Duration) {
	if t.cmd == nil {
		return
	}

	done := make(chan error, 1)
	go func() {
		done <- t.cmd.Wait()
	}()

	select {
	case <-time.After(shutdownTimeout):
		killProcessTree(t.cmd.Process.Pid)
	case err := <-done:
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func killProcessTree(pid int) {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	syscall.Kill(-pgid, syscall.SIGKILL)
}
