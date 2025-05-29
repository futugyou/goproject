package server

import (
	"bufio"
	"os"

	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/protocol"
)

type StdioServerTransport struct {
	*StreamServerTransport
}

func NewStdioServerTransport(serverName string, logger logging.Logger) *StdioServerTransport {
	if len(serverName) == 0 {
		serverName = "McpServer"
	}

	t := &StdioServerTransport{
		StreamServerTransport: NewStreamServerTransport(os.Stdin, bufio.NewWriter(os.Stdout), serverName, logger),
	}
	return t
}

func (t *StdioServerTransport) GetTransportKind() protocol.TransportKind {
	return protocol.TransportKindStdio
}
