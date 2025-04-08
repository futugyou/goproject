package server

type RequestContext[TParams any] struct {
	Params TParams
	Server IMcpServer
}
