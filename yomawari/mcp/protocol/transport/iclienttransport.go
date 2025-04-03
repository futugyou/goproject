package transport

import "context"

type IClientTransport interface {
	Connect(context.Context) (ITransport, error)
}
