package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	pb "github.com/goproject/tag-service/proto"
	"github.com/goproject/tag-service/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// protoc --go_out=plugins=grpc:. ./proto/*.proto
// grpcurl -plaintext  localhost:8001 proto.TagService.GetTagList
// protoc -IC:\Users\TerraformRs\Documents\GitHub\protobuf-3.13.0\include -I. -IC:\Users\TerraformRs\go -IC:\Users\TerraformRs\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.15.2\third_party\googleapis --go_out=plugins=grpc:.  --grpc-gateway_out=logtostderr=true:. ./proto/*.proto

var (
	port     string = "8001"
	grpcPort string
	httpPort string
)

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "8001", "grpc port")
	flag.StringVar(&httpPort, "hppt_port", "9001", "http port")
	flag.StringVar(&port, "port", "8001", "http&grpc port")
	flag.Parse()
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					grpcServer.ServeHTTP(w, r)
				} else {
					otherHandler.ServeHTTP(w, r)
				}
			}),
		&http2.Server{},
	)
}

func RunServer(port string) error {
	httpMux := runHttpServer()
	gatewayMux := runGrpcGatewayServer()
	grpcS := runGrpcServer()
	httpMux.Handle("/", gatewayMux)
	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}

func runHttpServer() *http.ServeMux {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
	return serverMux
}

func runGrpcServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}

func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	return gwmux
}
func RunTCPServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func RunHttpServer(port string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func main() {
	err := RunServer(port)
	if err != nil {
		log.Fatalf("run tcp server err: %v", err)
	}

	// l, err := RunTCPServer(port)
	// if err != nil {
	// 	log.Fatalf("run tcp server err: %v", err)
	// }
	// m := cmux.New(l)
	// grpcL := m.MatchWithWriters(
	// 	cmux.HTTP2MatchHeaderFieldPrefixSendSettings(
	// 		"content-type",
	// 		"application/grpc"))
	// httpL := m.Match(cmux.HTTP1Fast())

	// grpcS := RunGrpcServer()
	// httpS := RunHttpServer(port)

	// go grpcS.Serve(grpcL)
	// go httpS.Serve(httpL)
	// err = m.Serve()
	// if err != nil {
	// 	log.Fatalf("run serve err: %v", err)
	// }
}
