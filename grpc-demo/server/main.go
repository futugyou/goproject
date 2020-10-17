package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"

	pb "github.com/go-project/grpc-demo/proto"
	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "port")
	flag.Parse()
}

//protoc --go_out=plugins=grpc:. ./proto/*.proto
type GreeterServer struct{}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println("-------SayHello--------")
	return &pb.HelloReply{Message: "hello.world"}, nil
}

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	log.Println("-------SayList--------")
	for n := 0; n < 6; n++ {
		err := stream.Send(&pb.HelloReply{
			Message: "hello list",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	log.Println("-------SayRecord--------")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "say record"})
		}
		if err != nil {
			return nil
		}
		log.Printf("resp: %v", resp)
	}
}

func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	log.Println("-------SayRoute--------")
	n := 0
	for {
		err := stream.Send(&pb.HelloReply{Message: "say route= ="})
		if err != nil {
			return nil
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		n++
		log.Printf("resp: %v", resp)
	}
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":"+port)
	server.Serve(lis)
}
