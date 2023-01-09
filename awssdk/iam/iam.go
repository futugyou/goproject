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
		fmt.Println("UserName:", *user.UserName, "\tUserId:", *user.UserId, "\tPath:", *user.Path)
		fmt.Println("\tTags:", user.Tags)
		input := &iam.ListUserPoliciesInput{
			UserName: user.UserName,
		}
		output, err := svc.ListUserPolicies(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("\tPolicyNames:", output.PolicyNames)

		attachPolicyInput := &iam.ListAttachedUserPoliciesInput{
			UserName: input.UserName,
		}
		attachPolicyOutput, err := svc.ListAttachedUserPolicies(awsenv.EmptyContext, attachPolicyInput)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("\tAttachedPolicyName:")
		for _, policy := range attachPolicyOutput.AttachedPolicies {
			fmt.Println("\t- ", *policy.PolicyName)
		}

		fmt.Println()
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

		policyInput := &iam.ListGroupPoliciesInput{
			GroupName: group.GroupName,
		}

		policyOutput, err := svc.ListGroupPolicies(awsenv.EmptyContext, policyInput)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("\tPolicyNames:", policyOutput.PolicyNames)

		attachedPolicyInput := &iam.ListAttachedGroupPoliciesInput{
			GroupName: group.GroupName,
		}
		attachedPolicyOutput, err := svc.ListAttachedGroupPolicies(awsenv.EmptyContext, attachedPolicyInput)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("\tAttachedPolicyNames:")
		for _, policy := range attachedPolicyOutput.AttachedPolicies {
			fmt.Println("\t -", *policy.PolicyName)
		}
	}
}
