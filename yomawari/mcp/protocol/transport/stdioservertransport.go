package transport

import (
	"bufio"
	"os"

	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/server"
)

type StdioServerTransport struct {
	*StreamServerTransport
}

func NewStdioServerTransport(serverOptions *server.McpServerOptions, serverName string, logger logging.Logger) *StdioServerTransport {
	if len(serverName) == 0 && serverOptions != nil {
		serverName = serverOptions.ServerInfo.Name
	}
	t := &StdioServerTransport{
		StreamServerTransport: NewStreamServerTransport(os.Stdin, bufio.NewWriter(os.Stdout), serverName, logger),
	}
	return t
}
