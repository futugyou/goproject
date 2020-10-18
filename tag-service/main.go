package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/goproject/tag-service/proto"
	"github.com/goproject/tag-service/server"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// protoc --go_out=plugins=grpc:. ./proto/*.proto
// grpcurl -plaintext  localhost:8001 proto.TagService.GetTagList

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
	l, err := RunTCPServer(port)
	if err != nil {
		log.Fatalf("run tcp server err: %v", err)
	}
	m := cmux.New(l)
	grpcL := m.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldPrefixSendSettings(
			"content-type",
			"application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	grpcS := RunGrpcServer()
	httpS := RunHttpServer(port)

	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)
	err = m.Serve()
	if err != nil {
		log.Fatalf("run serve err: %v", err)
	}
}
