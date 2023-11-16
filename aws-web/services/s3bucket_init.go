package services

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

func (s *S3bucketService) GetBucketRegion(name *string, region string) string {
	svc := s3.NewFromConfig(awsenv.Cfg)

	input := s3.GetBucketLocationInput{
		Bucket: name,
	}
	output, err := svc.GetBucketLocation(awsenv.EmptyContext, &input)
	if err != nil {
		log.Println("GetBucketLocation error.")
		return region
	}

	return string(output.LocationConstraint)
}

func (s *S3bucketService) GetBucketPolicyStatus(name *string) bool {
	svc := s3.NewFromConfig(awsenv.Cfg)

	input := s3.GetBucketPolicyStatusInput{
		Bucket: name,
	}
	output, err := svc.GetBucketPolicyStatus(awsenv.EmptyContext, &input)
	if err != nil {
		log.Println("GetBucketPolicyStatus error.")
		return false
	}

	if output.PolicyStatus != nil {
		return output.PolicyStatus.IsPublic
	}

	return false
}

func (s *S3bucketService) GetBucketPolicy(name *string) string {
	svc := s3.NewFromConfig(awsenv.Cfg)

	input := s3.GetBucketPolicyInput{
		Bucket: name,
	}
	output, err := svc.GetBucketPolicy(awsenv.EmptyContext, &input)
	if err != nil {
		log.Println("GetBucketPolicy error.")
		return ""
	}

	if output.Policy != nil {
		return *output.Policy
	}

	return ""
}

func (s *S3bucketService) GetBucketPermissions(name *string) []string {
	svc := s3.NewFromConfig(awsenv.Cfg)
	input := s3.GetBucketAclInput{
		Bucket: name,
	}
	output, err := svc.GetBucketAcl(awsenv.EmptyContext, &input)
	if err != nil {
		log.Println("GetBucketAcl error.")
		return []string{}
	}

	result := make([]string, 0)
	for _, v := range output.Grants {
		result = append(result, string(v.Permission))
	}

	return result
}

func (s *S3bucketService) GetAllS3Bucket() ([]types.Bucket, error) {
	svc := s3.NewFromConfig(awsenv.Cfg)
	input := s3.ListBucketsInput{}

	output, err := svc.ListBuckets(awsenv.EmptyContext, &input)
	if err != nil {
		log.Println("get ListBuckets error.")
		return nil, err
	}

	return output.Buckets, nil
}
