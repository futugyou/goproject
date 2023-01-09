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
		fmt.Println("UserName:", *user.UserName, "\tUserId:", *user.UserId, "\tPath:", *user.Path, "\tTags:", user.Tags)
	}
}

func ListAccessKeys() {
	input := &iam.ListAccessKeysInput{}
	// this method will show currect user(env) key
	output, err := svc.ListAccessKeys(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, data := range output.AccessKeyMetadata {
		fmt.Println("UserName:", *data.UserName, "\tAccessKeyId:", *data.AccessKeyId)
	}
}
