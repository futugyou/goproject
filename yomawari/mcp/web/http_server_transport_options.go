package web

import (
	"context"
	"time"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/mcp/server"
)

type HttpContext struct {
	// TODO
}

type HttpServerTransportOptions struct {
	ConfigureSessionOptions func(context.Context, HttpContext, server.McpServerOptions) error
	RunSessionHandler       func(context.Context, HttpContext, server.IMcpServer) error
	Stateless               bool
	IdleTimeout             time.Duration
	MaxIdleSessionCount     int
	TimeProvider            core.TimeProvider
}

func NewHttpServerTransportOptions() *HttpServerTransportOptions {
	return &HttpServerTransportOptions{
		ConfigureSessionOptions: func(context.Context, HttpContext, server.McpServerOptions) error { return nil },
		RunSessionHandler:       func(context.Context, HttpContext, server.IMcpServer) error { return nil },
		Stateless:               false,
		IdleTimeout:             2 * time.Hour,
		MaxIdleSessionCount:     100_000,
		TimeProvider:            core.RealTimeProvider{},
	}
}
