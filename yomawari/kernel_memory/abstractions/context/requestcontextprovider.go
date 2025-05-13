package context

import "context"

type ContextKey string

const requestContextKey ContextKey = "RequestContext"

type RequestContextProvider struct {
	context *RequestContext
}

func (p *RequestContextProvider) GetContext(ctx context.Context) IContext {
	if p == nil {
		return nil
	}
	if ctx == nil {
		return p.context
	}
	val := ctx.Value(requestContextKey)
	if val == nil {
		return nil
	}
	return val.(*RequestContext)
}

func (p *RequestContextProvider) SetContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, requestContextKey, rc)
}

func (p *RequestContextProvider) InitContextArgs(ctx context.Context, args map[string]interface{}) IContextProvider {
	if p == nil {
		return nil
	}

	p.GetContext(ctx).InitArgs(args)
	return p
}

func (p *RequestContextProvider) InitContext(ctx context.Context, context IContext) IContextProvider {
	if p == nil {
		return nil
	}

	args := context.GetArgs()
	if args == nil {
		args = make(map[string]interface{})
	}

	p.GetContext(ctx).InitArgs(args)
	return p
}

func (p *RequestContextProvider) SetContextArgs(ctx context.Context, args map[string]interface{}) IContextProvider {
	if p == nil {
		return nil
	}

	p.GetContext(ctx).SetArgs(args)
	return p
}

func (p *RequestContextProvider) SetContextArg(ctx context.Context, key string, value interface{}) IContextProvider {
	if p == nil {
		return nil
	}

	p.GetContext(ctx).SetArg(key, value)
	return p
}
