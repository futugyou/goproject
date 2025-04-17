package transport

import "context"

type IClientTransport interface {
	GetName() string
	Connect(context.Context) (ITransport, error)
}
