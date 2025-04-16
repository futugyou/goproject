package transport

import (
	"bufio"
	"os"

	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

var McpServerDefaultImplementation types.Implementation = types.Implementation{
	Name:    "McpServer",
	Version: "1.0.0",
}

type StdioServerTransport struct {
	*StreamServerTransport
}

func NewStdioServerTransport(serverName string, logger logging.Logger) *StdioServerTransport {
	if len(serverName) == 0 {
		serverName = McpServerDefaultImplementation.Name
	}

	t := &StdioServerTransport{
		StreamServerTransport: NewStreamServerTransport(os.Stdin, bufio.NewWriter(os.Stdout), serverName, logger),
	}
	return t
}
