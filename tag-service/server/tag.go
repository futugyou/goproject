package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/goproject/tag-service/pkg/bapi"
	"github.com/goproject/tag-service/pkg/errcode"
	pb "github.com/goproject/tag-service/proto"
	"google.golang.org/grpc/metadata"
)

type TagServer struct{}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	dm, _ := metadata.FromIncomingContext(ctx)
	log.Printf("md: %+v", dm)
	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}
	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.Fail)
	}
	return &tagList, nil
}
