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

func ListGroups() {
	input := &iam.ListGroupsInput{}
	output, err := svc.ListGroups(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, group := range output.Groups {
		fmt.Print("GroupName:", *group.GroupName, "\tGroupId:", *group.GroupId)
		input := &iam.GetGroupInput{
			GroupName: group.GroupName,
		}
		output, err := svc.GetGroup(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("\tPath:", *output.Group.Path)
		for _, user := range output.Users {
			fmt.Println("\tUser:", *user.UserName)
		}
	}
}
