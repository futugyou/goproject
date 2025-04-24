package transport

import (
	"bufio"
	"os"

	"github.com/futugyou/yomawari/mcp/logging"
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
