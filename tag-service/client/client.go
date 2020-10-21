package main

import (
	"context"
	"log"

	"github.com/goproject/tag-service/global"
	"github.com/goproject/tag-service/internal/middleware"
	"github.com/goproject/tag-service/pkg/tracer"
	pb "github.com/goproject/tag-service/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func main() {
	auth := Auth{
		AppKey:    "appkey",
		AppSecret: "terraform",
	}
	ctx := context.Background()
	md := metadata.New(map[string]string{"hello": "goland", "one": "two"})
	// newCtx := metadata.AppendToOutgoingContext(ctx, "metadata-1", "go demo")
	newCtx := metadata.NewOutgoingContext(ctx, md)
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}
	clientConn, err := GetClientConn(newCtx, "localhost:8001", opts)
	if err != nil {
		log.Fatalf("err :%v", err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "golang"})
	if err != nil {
		log.Fatalf("tagserviceclient gettaglist err: %v", err)
	}
	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryContextTimeout(),
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithMax(2),
				grpc_retry.WithCodes(
					codes.Unknown,
					codes.Internal,
					codes.DeadlineExceeded,
				),
			),
			middleware.ClientTracing(),
			//otgrpc.OpenTracingClientInterceptor(global.Tracer),
		),
	))
	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(
			middleware.StreamContextTimeout(),
		),
	))

	// config := clientv3.Config{
	// 	Endpoints:   []string{"http://localhost:2379"},
	// 	DialTimeout: time.Second * 60,
	// }
	// cli, err := clientv3.New(config)
	// if err != nil {
	// 	return nil, err
	// }
	// r := &naming.GRPCResolver{Client: cli}
	// target := fmt.Sprintf("/etcdv3://goproject/grpc/%s", "tag-service")
	// opts = append(opts,
	// 	grpc.WithBalancer(grpc.RoundRobin(r)),
	// 	grpc, QithBlock(),
	// )

	return grpc.DialContext(ctx, target, opts...)
}

func init() {
	tacer, _, _ := tracer.NewJaegerTracer("tag_service", "127.0.0.1:6831")
	global.Tracer = tacer
}
