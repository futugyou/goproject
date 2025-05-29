package client

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/protocol"
)

type StdioClientSessionTransport struct {
	*StreamClientSessionTransport
	options *StdioClientTransportOptions
	cmd     *exec.Cmd
}

func NewStdioClientSessionTransport(options *StdioClientTransportOptions, cmd *exec.Cmd, endpointName string, logger logging.Logger) *StdioClientSessionTransport {
	return &StdioClientSessionTransport{
		StreamClientSessionTransport: NewStreamClientSessionTransport(cmd.Stdout, cmd.Stdin, endpointName, logger),
		options:                      options,
		cmd:                          cmd,
	}
}

func (t *StdioClientSessionTransport) SendMessage(ctx context.Context, message protocol.IJsonRpcMessage) error {
	if t.cmd.ProcessState != nil && t.cmd.ProcessState.Exited() {
		t.logger.TransportNotConnected(t.EndpointName)
		return fmt.Errorf("transport is not connected")
	}

	if err := t.cmd.Process.Signal(syscall.Signal(0)); err != nil {
		t.logger.TransportNotConnected(t.EndpointName)
		return fmt.Errorf("transport is not connected: %v", err)
	}

	return t.StreamClientSessionTransport.SendMessage(ctx, message)
}

func (t *StdioClientSessionTransport) Close() error {
	if t.cmd == nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- t.cmd.Wait()
	}()

	select {
	case <-time.After(t.options.ShutdownTimeout):
		killProcessTree(t.cmd.Process.Pid)
	case err := <-done:
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return t.StreamClientSessionTransport.Close()
}

func (t *StdioClientSessionTransport) GetTransportKind() protocol.TransportKind {
	return protocol.TransportKindStdio
}
