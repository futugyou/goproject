package main

import (
	"context"
	"io"
	"log"

	pb "github.com/go-project/grpc-demo/proto"
	"google.golang.org/grpc"
)

var port string = "8000"

func main() {
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	request := &pb.HelloRequest{
		Name: "toy",
	}

	log.Println("-------SayHello--------")
	err := SayHello(client)
	log.Println("-------SayList--------")
	err = SayList(client, request)
	log.Println("-------SayRecord--------")
	err = SayRecord(client, request)
	log.Println("-------SayRoute--------")
	err = SayRoute(client, request)
	if err != nil {
		log.Fatalf("sayhello err: %v", err)
	}
}

func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRoute(context.Background())
	if err != nil {
		return nil
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
	err = stream.CloseSend()
	if err != nil {
		return err
	}
	return nil
}

func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRecord(context.Background())
	if err != nil {
		return nil
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp err: %v", resp)
	return nil
}

func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayList(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
	return nil
}

func SayHello(client pb.GreeterClient) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
		Name: "toy",
	})
	if err != nil {
		return err
	}
	log.Printf("client. sayhello resp: %s", resp.Message)
	return nil
}
