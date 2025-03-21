package context

import "context"

type IContextProvider interface {
	GetContext(ctx context.Context) IContext
	SetContext(ctx context.Context, rc *RequestContext) context.Context
	InitContextArgs(ctx context.Context, args map[string]interface{}) IContextProvider
	InitContext(ctx context.Context, context IContext) IContextProvider
	SetContextArgs(ctx context.Context, args map[string]interface{}) IContextProvider
	SetContextArg(ctx context.Context, key string, value interface{}) IContextProvider
}

var _ IContextProvider = &RequestContextProvider{}
