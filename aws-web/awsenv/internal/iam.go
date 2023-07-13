package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

var (
	svc             *iam.Client
	internalContext context.Context
)

func CompleteUser(cfg aws.Config, groupName, userName, groupPolicyArn string) {
	svc = iam.NewFromConfig(cfg)
	internalContext = context.Background()

	// 1.create iam group
	input := &iam.CreateGroupInput{
		GroupName: aws.String(groupName),
	}
	_, err := svc.CreateGroup(internalContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. attach policy to group
	policyInput := &iam.AttachGroupPolicyInput{
		GroupName: aws.String(groupName),
		PolicyArn: aws.String(groupPolicyArn),
	}
	_, err = svc.AttachGroupPolicy(internalContext, policyInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3.create iam user
	userInput := &iam.CreateUserInput{
		UserName: aws.String(userName),
	}
	_, err = svc.CreateUser(internalContext, userInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4.create iam access key
	keyInput := &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	}
	keyOutput, err := svc.CreateAccessKey(internalContext, keyInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 5. add user to group
	userGroupInput := &iam.AddUserToGroupInput{
		UserName:  aws.String(userName),
		GroupName: aws.String(groupName),
	}
	_, err = svc.AddUserToGroup(internalContext, userGroupInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Setenv("AWS_ACCESS_KEY_ID", *keyOutput.AccessKey.AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", *keyOutput.AccessKey.SecretAccessKey)
}

func DeleteUser(cfg aws.Config, groupName, userName, groupPolicyArn string) {
	// 1. Detach Group Policy
	policyInput := &iam.DetachGroupPolicyInput{
		GroupName: aws.String(groupName),
		PolicyArn: aws.String(groupPolicyArn),
	}
	_, err := svc.DetachGroupPolicy(internalContext, policyInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Delete Access Key
	keyInput := &iam.DeleteAccessKeyInput{
		UserName:    aws.String(userName),
		AccessKeyId: aws.String(os.Getenv("AWS_ACCESS_KEY_ID")),
	}
	_, err = svc.DeleteAccessKey(internalContext, keyInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3. Remove User From Group
	userGroupInput := &iam.RemoveUserFromGroupInput{
		UserName:  aws.String(userName),
		GroupName: aws.String(groupName),
	}
	_, err = svc.RemoveUserFromGroup(internalContext, userGroupInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4. Delete User
	userInput := &iam.DeleteUserInput{
		UserName: aws.String(userName),
	}
	_, err = svc.DeleteUser(internalContext, userInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 5. Delete Group
	input := &iam.DeleteGroupInput{
		GroupName: aws.String(groupName),
	}
	_, err = svc.DeleteGroup(internalContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
}
