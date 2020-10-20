package main

import (
	"context"
	"log"

	"github.com/goproject/tag-service/internal/middleware"
	pb "github.com/goproject/tag-service/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	md := metadata.New(map[string]string{"hello": "goland", "one": "two"})
	// newCtx := metadata.AppendToOutgoingContext(ctx, "metadata-1", "go demo")
	newCtx := metadata.NewOutgoingContext(ctx, md)
	clientConn, err := GetClientConn(newCtx, "localhost:8001", nil)
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
		),
	))
	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(
			middleware.StreamContextTimeout(),
		),
	))

	return grpc.DialContext(ctx, target, opts...)
}
