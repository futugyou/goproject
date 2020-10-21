package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/goproject/blog-service/global"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var ctx context.Context
		// span := opentracing.SpanFromContext(c.Request.Context())
		// if span != nil {
		// 	span, ctx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(),
		// 		global.Tracer,
		// 		c.Request.URL.Path,
		// 		opentracing.ChildOf(span.Context()))
		// } else {
		// 	span, ctx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(),
		// 		global.Tracer,
		// 		c.Request.URL.Path)
		// }
		// defer span.Finish()

		var newCtx context.Context
		var span opentracing.Span
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), global.Tracer, c.Request.URL.Path)
		} else {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				global.Tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		defer span.Finish()

		var tracid string
		var spanid string
		var spanContext = span.Context()
		switch spanContext.(type) {
		case jaeger.SpanContext:
			tracid = spanContext.(jaeger.SpanContext).TraceID().String()
			spanid = spanContext.(jaeger.SpanContext).SpanID().String()
		}
		c.Set("X-Trace-ID", tracid)
		c.Set("X-Span-ID", spanid)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
