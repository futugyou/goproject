package extensions

import "context"

type ctxKey string

const JWTKey ctxKey = "internaljwt"

func WithJWT(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, JWTKey, token)
}

func JWTFrom(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(JWTKey).(string)
	return token, ok
}
