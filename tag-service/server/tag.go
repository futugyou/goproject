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

type TagServer struct {
	auth *Auth
}

type Auth struct{}

func (a *Auth) GetAppKey() string {
	return "appkey"
}
func (a *Auth) GetAppSecret() string {
	return "terraform"
}

func (a *Auth) check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("md: %+v", md)

	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
	return nil
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	// if err := t.auth.check(ctx); err != nil {
	// 	return nil, err
	// }
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
