package iot

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc             *iot.Client
	svciotdataplane *iotdataplane.Client
)

func init() {
	svc = iot.NewFromConfig(awsenv.Cfg)
	svciotdataplane = iotdataplane.NewFromConfig(awsenv.Cfg)
}

func ListJobs() {
	var nextToken *string = nil
	for {
		input := &iot.ListJobsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListJobs(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, job := range result.Jobs {
			log.Println(*job.JobId, job.Status, *job.ThingGroupId, job.TargetSelection)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThings() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListThings(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, thing := range result.Things {
			log.Println(*thing.ThingArn, *thing.ThingName, thing.ThingTypeName, thing.Version)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThingTypes() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingTypesInput{
			NextToken: nextToken,
		}

		result, err := svc.ListThingTypes(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, ty := range result.ThingTypes {
			log.Println(*ty.ThingTypeName, *ty.ThingTypeArn)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThingGroups() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingGroupsInput{
			NextToken: nextToken,
			Recursive: aws.Bool(true),
		}

		result, err := svc.ListThingGroups(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, group := range result.ThingGroups {
			log.Println(*group.GroupArn, *group.GroupName)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThingRegistrationTasks() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingRegistrationTasksInput{
			NextToken: nextToken,
		}

		result, err := svc.ListThingRegistrationTasks(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, taskid := range result.TaskIds {
			log.Println(taskid)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListTopicRuleDestinations() {
	var nextToken *string = nil
	for {
		input := &iot.ListTopicRuleDestinationsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListTopicRuleDestinations(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, summ := range result.DestinationSummaries {
			log.Println(*summ.Arn, summ.Status, *summ.StatusReason, summ.HttpUrlSummary, summ.VpcDestinationSummary)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListTopicRules() {
	var nextToken *string = nil
	for {
		input := &iot.ListTopicRulesInput{
			NextToken: nextToken,
		}

		result, err := svc.ListTopicRules(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, rule := range result.Rules {
			log.Println(*rule.RuleArn, *rule.RuleDisabled, *rule.RuleName, *rule.TopicPattern, *rule.CreatedAt)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListActiveViolations() {
	var nextToken *string = nil
	for {
		input := &iot.ListActiveViolationsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListActiveViolations(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, item := range result.ActiveViolations {
			log.Println(*item.Behavior)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListStreams() {
	var nextToken *string = nil
	for {
		input := &iot.ListStreamsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListStreams(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, item := range result.Streams {
			log.Println(*item.StreamArn, *item.StreamId)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func DescribeThing() {
	input := &iot.DescribeThingInput{
		ThingName: aws.String(awsenv.IotThingName),
	}

	result, err := svc.DescribeThing(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if result.BillingGroupName != nil {
		log.Println("BillingGroupName:\t", *result.BillingGroupName)
	}

	if result.DefaultClientId != nil {
		log.Println("DefaultClientId:\t", *result.DefaultClientId)
	}

	if result.ThingArn != nil {
		log.Println("ThingArn:\t", *result.ThingArn)
	}

	if result.ThingId != nil {
		log.Println("ThingId:\t", *result.ThingId)
	}

	if result.ThingName != nil {
		log.Println("ThingName:\t", *result.ThingName)
	}

	if result.ThingTypeName != nil {
		log.Println("ThingTypeName:\t", *result.ThingTypeName)
	}

	for key, value := range result.Attributes {
		log.Println(key, "  ", value)
	}
}

func DescribeThingGroup() {
	input := &iot.DescribeThingGroupInput{
		ThingGroupName: aws.String(awsenv.IotThingGroupName),
	}

	result, err := svc.DescribeThingGroup(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if result.IndexName != nil {
		log.Println("IndexName:\t", *result.IndexName)
	}

	if result.QueryString != nil {
		log.Println("QueryString:\t", *result.QueryString)
	}

	if result.QueryVersion != nil {
		log.Println("QueryVersion:\t", *result.QueryVersion)
	}

	log.Println("Status:\t", result.Status)

	if result.ThingGroupArn != nil {
		log.Println("ThingGroupArn:\t", *result.ThingGroupArn)
	}

	if result.ThingGroupId != nil {
		log.Println("ThingGroupId:\t", *result.ThingGroupId)
	}

	if result.ThingGroupName != nil {
		log.Println("ThingGroupName:\t", *result.ThingGroupName)
	}

	log.Println("ThingGroupMetadata:\t", result.ThingGroupMetadata)
	log.Println("ThingGroupProperties:\t", result.ThingGroupProperties)
}

func GetRegistrationCode() {
	input := &iot.GetRegistrationCodeInput{}

	result, err := svc.GetRegistrationCode(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if result.RegistrationCode != nil {
		log.Println("RegistrationCode:\t", *result.RegistrationCode)
	}
}

func DescribeEndpoint() {
	endpointtypes := []string{"iot:Data", "iot:Data-ATS", "iot:CredentialProvider", "iot:Jobs"}
	for _, t := range endpointtypes {
		input := &iot.DescribeEndpointInput{
			EndpointType: aws.String(t),
		}

		result, err := svc.DescribeEndpoint(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if result.EndpointAddress != nil {
			log.Println(t, " :\t", *result.EndpointAddress)
		}
	}
}

func ListBillingGroups() {
	var nextToken *string = nil
	for {
		input := &iot.ListBillingGroupsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListBillingGroups(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, item := range result.BillingGroups {
			log.Println(*item.GroupArn, *item.GroupName)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListAuthorizers() {
	var nextToken *string = nil
	for {
		input := &iot.ListAuthorizersInput{
			Marker: nextToken,
		}

		result, err := svc.ListAuthorizers(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextMarker

		for _, item := range result.Authorizers {
			log.Println(*item.AuthorizerArn, *item.AuthorizerName)
		}
		if result.NextMarker == nil {
			return
		}
	}
}

func ListCACertificates() {
	var nextToken *string = nil
	for {
		input := &iot.ListCACertificatesInput{
			Marker: nextToken,
		}

		result, err := svc.ListCACertificates(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextMarker

		for _, item := range result.Certificates {
			log.Println(*item.CertificateArn, *item.CertificateId, item.Status)
		}
		if result.NextMarker == nil {
			return
		}
	}
}

func ListCertificates() {
	var nextToken *string = nil
	for {
		input := &iot.ListCertificatesInput{
			Marker: nextToken,
		}

		result, err := svc.ListCertificates(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextMarker

		for _, item := range result.Certificates {
			log.Println(*item.CertificateArn, *item.CertificateId, item.Status)
		}
		if result.NextMarker == nil {
			return
		}
	}
}

func ListPolicies() {
	var nextToken *string = nil
	for {
		input := &iot.ListPoliciesInput{
			Marker: nextToken,
		}

		result, err := svc.ListPolicies(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextMarker

		for _, item := range result.Policies {
			log.Println(*item.PolicyArn, *item.PolicyName)
			GetPolicy(*item.PolicyName)
			log.Println()
		}
		if result.NextMarker == nil {
			return
		}
	}
}

func GetPolicy(name string) {
	input := &iot.GetPolicyInput{
		PolicyName: aws.String(name),
	}

	result, err := svc.GetPolicy(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if result.CreationDate != nil {
		log.Println("CreationDate:\t", *result.CreationDate)
	}
	if result.DefaultVersionId != nil {
		log.Println("DefaultVersionId:\t", *result.DefaultVersionId)
	}
	if result.GenerationId != nil {
		log.Println("GenerationId:\t", *result.GenerationId)
	}
	if result.LastModifiedDate != nil {
		log.Println("LastModifiedDate:\t", *result.LastModifiedDate)
	}
	if result.PolicyArn != nil {
		log.Println("PolicyArn:\t", *result.PolicyArn)
	}
	if result.PolicyDocument != nil {
		log.Println("PolicyDocument:\t", *result.PolicyDocument)
	}
	if result.PolicyName != nil {
		log.Println("PolicyName:\t", *result.PolicyName)
	}
}

func ListRetainedMessages() {
	var nextToken *string = nil
	for {
		input := &iotdataplane.ListRetainedMessagesInput{
			NextToken: nextToken,
		}

		result, err := svciotdataplane.ListRetainedMessages(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, item := range result.RetainedTopics {
			log.Println(*item.Topic, item.Qos, item.PayloadSize)
			GetRetainedMessage(*item.Topic)
			log.Println()
		}
		if result.NextToken == nil {
			return
		}
	}
}

func GetRetainedMessage(name string) {
	input := &iotdataplane.GetRetainedMessageInput{
		Topic: aws.String(name),
	}

	result, err := svciotdataplane.GetRetainedMessage(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Payload:\t", string(result.Payload))
}

func ListNamedShadowsForThing() {
	var nextToken *string = nil
	for {
		input := &iotdataplane.ListNamedShadowsForThingInput{
			NextToken: nextToken,
			ThingName: aws.String(awsenv.IotThingName),
		}

		result, err := svciotdataplane.ListNamedShadowsForThing(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken
		log.Println("Results:\t", strings.Join(result.Results, ","))
		if result.NextToken == nil {
			return
		}
	}
}

func GetThingShadow() {
	input := &iotdataplane.GetThingShadowInput{
		ThingName: aws.String(awsenv.IotThingName),
	}

	result, err := svciotdataplane.GetThingShadow(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Payload:\t", string(result.Payload))
}
