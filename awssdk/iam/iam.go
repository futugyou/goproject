package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	input := &iam.ListAccessKeysInput{
		// UserName: aws.String(awsenv.UserName),
	}
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
			fmt.Println("\t -", *policy.PolicyName, *policy.PolicyArn)
		}
	}
}

func ListAccountAliases() {
	input := &iam.ListAccountAliasesInput{}
	output, err := svc.ListAccountAliases(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("AccountAliases:", output.AccountAliases)
}

func CreateAccountAlias() {
	input := &iam.CreateAccountAliasInput{
		AccountAlias: aws.String("jenkins-account"),
	}
	output, err := svc.CreateAccountAlias(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ResultMetadata:", output.ResultMetadata)
}

func DeleteAccountAlias() {
	input := &iam.DeleteAccountAliasInput{
		AccountAlias: aws.String("jenkins-account"),
	}
	output, err := svc.DeleteAccountAlias(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ResultMetadata:", output.ResultMetadata)
}

func ListInstanceProfiles() {
	input := &iam.ListInstanceProfilesInput{}
	output, err := svc.ListInstanceProfiles(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, profile := range output.InstanceProfiles {
		fmt.Println("Name:", *profile.InstanceProfileName, "\tId:", *profile.InstanceProfileId, "\tPath:", *profile.Path)
		fmt.Println("\tTags:")
		for _, tag := range profile.Tags {
			fmt.Println("\t -", *tag.Key, *tag.Value)
		}
		fmt.Println("\tRoles:")
		for _, role := range profile.Roles {
			fmt.Println("\t -", *role.RoleName, *role.AssumeRolePolicyDocument)
		}
	}
}

func ListPolicies() {
	input := &iam.ListPoliciesInput{}
	output, err := svc.ListPolicies(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, policy := range output.Policies {
		fmt.Println("Name:", *policy.PolicyName)
	}
}

func ListRoles() {
	input := &iam.ListRolesInput{}
	output, err := svc.ListRoles(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, role := range output.Roles {
		fmt.Println("RoleName:", *role.RoleName, "\tPath:", *role.Path)
		input := &iam.ListRolePoliciesInput{
			RoleName: role.RoleName,
		}
		output, err := svc.ListRolePolicies(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("\tPolicyNames:", output.PolicyNames)
	}
}

func CreateGroup() {
	input := &iam.CreateGroupInput{
		GroupName: aws.String(awsenv.GroupName),
	}
	output, err := svc.CreateGroup(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Name:", *output.Group.GroupName, "\tId:", *output.Group.GroupId, "\tPath:", *output.Group.Path)
}

func DeleteGroup() {
	input := &iam.DeleteGroupInput{
		GroupName: aws.String(awsenv.GroupName),
	}
	output, err := svc.DeleteGroup(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ResultMetadata:", output.ResultMetadata)
}

func CreateUser() {
	input := &iam.CreateUserInput{
		UserName: aws.String(awsenv.UserName),
	}
	output, err := svc.CreateUser(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("UserName:", *output.User.UserName, "\tUserId:", *output.User.UserId, "\tPath:", *output.User.Path)

	keyInput := &iam.CreateAccessKeyInput{
		UserName: output.User.UserName,
	}
	keyOutput, err := svc.CreateAccessKey(awsenv.EmptyContext, keyInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("\tAccessKeyId:", *keyOutput.AccessKey.AccessKeyId, "\tSecretAccessKey:", *keyOutput.AccessKey.SecretAccessKey)
}

func DeleteUser() {
	// must delete referce resouce before delete user
	keyInput := &iam.DeleteAccessKeyInput{
		UserName:    aws.String(awsenv.UserName),
		AccessKeyId: aws.String("AKIAVSXAPJPTNZK77PLI"),
	}
	_, err := svc.DeleteAccessKey(awsenv.EmptyContext, keyInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	input := &iam.DeleteUserInput{
		UserName: aws.String(awsenv.UserName),
	}
	output, err := svc.DeleteUser(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ResultMetadata:", output.ResultMetadata)
}

func GetAccountAuthorizationDetails() {
	input := &iam.GetAccountAuthorizationDetailsInput{}
	output, err := svc.GetAccountAuthorizationDetails(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GroupDetailList:")
	for _, group := range output.GroupDetailList {
		fmt.Println("\tGroupName:", *group.GroupName)
		for _, policy := range group.GroupPolicyList {
			fmt.Println("\t - ", *policy.PolicyName)
		}
	}

	fmt.Println()
	fmt.Println("Policies:")
	for _, policy := range output.Policies {
		fmt.Println("\tPolicyName:", *policy.PolicyName)
	}

	fmt.Println()
	fmt.Println("RoleDetailList:")
	for _, role := range output.RoleDetailList {
		fmt.Println("\tRoleName:", *role.RoleName)
		for _, policy := range role.RolePolicyList {
			fmt.Println("\t\tPolicyName:", *policy.PolicyName)
		}
	}

	fmt.Println()
	fmt.Println("UserDetailList:")
	for _, user := range output.UserDetailList {
		fmt.Println("\tUserName:", *user.UserName, "\tGroupList:", user.GroupList)
		for _, policy := range user.UserPolicyList {
			fmt.Println("\t\tPolicyName:", *policy.PolicyName)
		}
	}
}
