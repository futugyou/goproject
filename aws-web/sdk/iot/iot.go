package iot

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
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
	endpointtypes := []string{"iot:Data", "iot:Data-ATS", "iot:CredentialProvider", "iot:Jobs", "iot:DeviceAdvisor", "iot:Data-Beta"}
	for _, t := range endpointtypes {
		input := &iot.DescribeEndpointInput{
			EndpointType: aws.String(t),
		}

		result, err := svc.DescribeEndpoint(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			continue
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
			DescribeCACertificate(*item.CertificateId)
			log.Println()
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
			DescribeCertificate(*item.CertificateId)
			log.Println()
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

func ListDomainConfigurations() {
	var nextToken *string = nil
	for {
		input := &iot.ListDomainConfigurationsInput{
			Marker: nextToken,
		}

		result, err := svc.ListDomainConfigurations(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextMarker
		for _, item := range result.DomainConfigurations {
			log.Println(*item.DomainConfigurationArn, *item.DomainConfigurationName, item.ServiceType)
			DescribeDomainConfiguration(*item.DomainConfigurationName)
			log.Println()
		}
		if result.NextMarker == nil {
			return
		}
	}
}

func DescribeDomainConfiguration(name string) {
	input := &iot.DescribeDomainConfigurationInput{
		DomainConfigurationName: aws.String(name),
	}

	result, err := svc.DescribeDomainConfiguration(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if result.AuthorizerConfig != nil {
		if result.AuthorizerConfig.AllowAuthorizerOverride != nil {
			log.Println("AllowAuthorizerOverride:\t", *result.AuthorizerConfig.AllowAuthorizerOverride)
		}
		if result.AuthorizerConfig.DefaultAuthorizerName != nil {
			log.Println("DefaultAuthorizerName:\t", *result.AuthorizerConfig.DefaultAuthorizerName)
		}
	}

	if result.DomainConfigurationArn != nil {
		log.Println("DomainConfigurationArn:\t", *result.DomainConfigurationArn)
	}

	if result.DomainConfigurationName != nil {
		log.Println("DomainConfigurationName:\t", *result.DomainConfigurationName)
	}

	if result.DomainName != nil {
		log.Println("DomainName:\t", *result.DomainName)
	}

	if result.LastStatusChangeDate != nil {
		log.Println("LastStatusChangeDate:\t", *result.LastStatusChangeDate)
	}

	if result.TlsConfig != nil {
		if result.TlsConfig.SecurityPolicy != nil {
			log.Println("SecurityPolicy:\t", *result.TlsConfig.SecurityPolicy)
		}
	}

	log.Println("ServiceType:\t", result.ServiceType)

	log.Println("DomainConfigurationStatus:\t", result.DomainConfigurationStatus)

	log.Println("DomainType:\t", result.DomainType)

	for _, item := range result.ServerCertificates {
		if item.ServerCertificateArn != nil {
			log.Println("ServerCertificateArn:\t", *item.ServerCertificateArn)
		}
		if item.ServerCertificateStatusDetail != nil {
			log.Println("ServerCertificateStatusDetail:\t", *item.ServerCertificateStatusDetail)
		}
		log.Println("ServerCertificateStatus:\t", item.ServerCertificateStatus)
	}

}

func DescribeCACertificate(name string) {
	input := &iot.DescribeCACertificateInput{
		CertificateId: aws.String(name),
	}

	result, err := svc.DescribeCACertificate(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	showCACertificateDescription(result.CertificateDescription)

	if result.RegistrationConfig != nil {
		if result.RegistrationConfig.RoleArn != nil {
			log.Println("RoleArn:\t", *result.RegistrationConfig.RoleArn)
		}
		if result.RegistrationConfig.TemplateBody != nil {
			log.Println("TemplateBody:\t", *result.RegistrationConfig.TemplateBody)
		}
		if result.RegistrationConfig.TemplateName != nil {
			log.Println("TemplateName:\t", *result.RegistrationConfig.TemplateName)
		}
	}
}

func showCACertificateDescription(CertificateDescription *types.CACertificateDescription) {
	if CertificateDescription != nil {
		if CertificateDescription.CertificateArn != nil {
			log.Println("CertificateArn:\t", *CertificateDescription.CertificateArn)
		}
		if CertificateDescription.CertificateId != nil {
			log.Println("CertificateId:\t", *CertificateDescription.CertificateId)
		}
		if CertificateDescription.CertificatePem != nil {
			log.Println("CertificatePem:\t", *CertificateDescription.CertificatePem)
		}
		if CertificateDescription.CreationDate != nil {
			log.Println("CreationDate:\t", *CertificateDescription.CreationDate)
		}
		if CertificateDescription.CustomerVersion != nil {
			log.Println("CustomerVersion:\t", *CertificateDescription.CustomerVersion)
		}
		if CertificateDescription.GenerationId != nil {
			log.Println("GenerationId:\t", *CertificateDescription.GenerationId)
		}
		if CertificateDescription.LastModifiedDate != nil {
			log.Println("LastModifiedDate:\t", *CertificateDescription.LastModifiedDate)
		}
		if CertificateDescription.OwnedBy != nil {
			log.Println("OwnedBy:\t", *CertificateDescription.OwnedBy)
		}
		if CertificateDescription.Validity != nil {
			if CertificateDescription.Validity.NotAfter != nil {
				log.Println("NotAfter:\t", *CertificateDescription.Validity.NotAfter)
			}
			if CertificateDescription.Validity.NotBefore != nil {
				log.Println("NotBefore:\t", *CertificateDescription.Validity.NotBefore)
			}
		}

		log.Println("AutoRegistrationStatus:\t", CertificateDescription.AutoRegistrationStatus)
		log.Println("CertificateMode:\t", CertificateDescription.CertificateMode)
		log.Println("Status:\t", CertificateDescription.Status)
	}
}

func DescribeCertificate(name string) {
	input := &iot.DescribeCertificateInput{
		CertificateId: aws.String(name),
	}

	result, err := svc.DescribeCertificate(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	showCertificateDescription(result.CertificateDescription)
}

func showCertificateDescription(CertificateDescription *types.CertificateDescription) {
	if CertificateDescription != nil {
		if CertificateDescription.CertificateArn != nil {
			log.Println("CertificateArn:\t", *CertificateDescription.CertificateArn)
		}
		if CertificateDescription.CertificateId != nil {
			log.Println("CertificateId:\t", *CertificateDescription.CertificateId)
		}
		if CertificateDescription.CertificatePem != nil {
			log.Println("CertificatePem:\t", *CertificateDescription.CertificatePem)
		}
		if CertificateDescription.CreationDate != nil {
			log.Println("CreationDate:\t", *CertificateDescription.CreationDate)
		}
		if CertificateDescription.CustomerVersion != nil {
			log.Println("CustomerVersion:\t", *CertificateDescription.CustomerVersion)
		}
		if CertificateDescription.GenerationId != nil {
			log.Println("GenerationId:\t", *CertificateDescription.GenerationId)
		}
		if CertificateDescription.LastModifiedDate != nil {
			log.Println("LastModifiedDate:\t", *CertificateDescription.LastModifiedDate)
		}
		if CertificateDescription.OwnedBy != nil {
			log.Println("OwnedBy:\t", *CertificateDescription.OwnedBy)
		}
		if CertificateDescription.Validity != nil {
			if CertificateDescription.Validity.NotAfter != nil {
				log.Println("NotAfter:\t", *CertificateDescription.Validity.NotAfter)
			}
			if CertificateDescription.Validity.NotBefore != nil {
				log.Println("NotBefore:\t", *CertificateDescription.Validity.NotBefore)
			}
		}
		if CertificateDescription.PreviousOwnedBy != nil {
			log.Println("PreviousOwnedBy:\t", *CertificateDescription.PreviousOwnedBy)
		}

		log.Println("CertificateMode:\t", CertificateDescription.CertificateMode)
		log.Println("Status:\t", CertificateDescription.Status)
		if CertificateDescription.TransferData != nil {
			if CertificateDescription.TransferData.AcceptDate != nil {
				log.Println("AcceptDate:\t", *CertificateDescription.TransferData.AcceptDate)
			}
			if CertificateDescription.TransferData.RejectDate != nil {
				log.Println("RejectDate:\t", *CertificateDescription.TransferData.RejectDate)
			}
			if CertificateDescription.TransferData.RejectReason != nil {
				log.Println("RejectReason:\t", *CertificateDescription.TransferData.RejectReason)
			}
			if CertificateDescription.TransferData.TransferDate != nil {
				log.Println("TransferDate:\t", *CertificateDescription.TransferData.TransferDate)
			}
			if CertificateDescription.TransferData.TransferMessage != nil {
				log.Println("TransferMessage:\t", *CertificateDescription.TransferData.TransferMessage)
			}
		}
	}
}
