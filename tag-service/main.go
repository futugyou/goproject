package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/goproject/tag-service/global"
	"github.com/goproject/tag-service/internal/middleware"
	"github.com/goproject/tag-service/pkg/swagger"
	"github.com/goproject/tag-service/pkg/tracer"
	pb "github.com/goproject/tag-service/proto"
	"github.com/goproject/tag-service/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// protoc --go_out=plugins=grpc:. ./proto/*.proto
// grpcurl -plaintext  localhost:8001 proto.TagService.GetTagList
// protoc -IC:\Users\TerraformRs\Documents\GitHub\protobuf-3.13.0\include -I. -IC:\Users\TerraformRs\go -IC:\Users\TerraformRs\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.15.2\third_party\googleapis --go_out=plugins=grpc:. --swagger_out=logtostderr=true:.  --grpc-gateway_out=logtostderr=true:. ./proto/*.proto
// go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/...

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
	tacer, _, _ := tracer.NewJaegerTracer("tag_service", "127.0.0.1:6831")
	global.Tracer = tacer
}

type httpError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func grpcGatewayError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}
	httpError := httpError{Code: int32(s.Code()), Message: s.Message()}
	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*pb.Error); ok {
			httpError.Code = v.Code
			httpError.Message = v.Message
		}
	}

	resp, _ := json.Marshal(httpError)
	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	_, _ = w.Write(resp)
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

	// etcdClient, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{"http://localhost:2379"},
	// 	DialTimeout: time.Second * 60,
	// })
	// if err != nil {
	// 	return err
	// }
	// defer etcdClient.Close()

	// target := fmt.Sprintf("/etcdv3://goproject/grpc/%s", "tag-service")
	// grpcproxy.Register(etcdClient, target, ":"+port, 60)

	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}

func runHttpServer() *http.ServeMux {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serverMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	serverMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}
		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("proto", p)
		http.ServeFile(w, r, p)
	})
	return serverMux
}

func runGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// middleware.AccessLog,
			// middleware.ErrorLog,
			// middleware.Recovery,
			middleware.ServerTracing,
		//otgrpc.OpenTracingServerInterceptor(global.Tracer),
		)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}
func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	runtime.HTTPError = grpcGatewayError
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
