package middleware

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/goproject/tag-service/global"
	"github.com/goproject/tag-service/pkg/errcode"
	"github.com/goproject/tag-service/pkg/metatext"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("Hello one")
	resp, err := handler(ctx, req)
	log.Println("bye one!")
	return resp, err
}

func WorldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("Hello two")
	resp, err := handler(ctx, req)
	log.Println("bye two!")
	return resp, err
}

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestLog := "access request log: method: %s, begin_time: %d, request: %v"
	beginTime := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, beginTime, req)
	resp, err := handler(ctx, req)
	responseLog := "access response log:method: %s,begin_time: %d, end_time: %d,response: %v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
	return resp, err
}

func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		errLog := "error log:method: %s, code: %v, message: %v, details: %v"
		s := errcode.FromError(err)
		log.Printf(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
	}
	return resp, err
}

func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			recoveryLog := "recover log: method: %s, message: %v, stack: %s"
			log.Printf(recoveryLog, info.FullMethod, e, string(debug.Stack()[:]))
		}
	}()
	return handler(ctx, req)
}

func ServerTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	parentSpanContext, _ := global.Tracer.Extract(opentracing.TextMap, metatext.MetadataTextMap{md})
	spanOpts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: string(ext.Component), Value: "grpc"},
		ext.SpanKindRPCServer,
		ext.RPCServerOption(parentSpanContext),
	}

	span := global.Tracer.StartSpan(info.FullMethod, spanOpts...)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return handler(ctx, req)
}
