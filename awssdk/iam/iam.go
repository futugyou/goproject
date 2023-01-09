package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *iam.Client
)

func init() {
	svc = iam.NewFromConfig(awsenv.Cfg)
}

func ListUsers() {
	input := &iam.ListUsersInput{}

	output, err := svc.ListUsers(awsenv.EmptyContext, input)

	if err != nil {
		fmt.Println(err)
		return
	}

	// this will be nil
	//fmt.Println("Marker:", *output.Marker)
	for _, user := range output.Users {
		fmt.Println(*user.UserName, *user.UserId, *user.Path, user.Tags)
	}
}
