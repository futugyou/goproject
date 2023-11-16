package services

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

func (s *S3bucketService) InitData() {
	log.Println("s3 sync start..")
	ctx := context.Background()

	accountService := NewAccountService()
	accounts := accountService.GetAllAccounts()
	
	entities := make([]entity.S3bucketEntity, 0)
	for _, account := range accounts {
		awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
		buckets, err := s.GetAllS3Bucket()
		if err != nil {
			continue
		}

		for _, bucket := range buckets {
			b := entity.S3bucketEntity{
				Id:           account.Id + *bucket.Name,
				Name:         *bucket.Name,
				Region:       s.GetBucketRegion(bucket.Name, account.Region),
				IsPublic:     s.GetBucketPolicyStatus(bucket.Name),
				Policy:       s.GetBucketPolicy(bucket.Name),
				Permissions:  s.GetBucketPermissions(bucket.Name),
				CreationDate: *bucket.CreationDate,
			}
			entities = append(entities, b)
		}
	}

	if len(entities) > 0 {
		s.repository.DeleteAll(ctx)
		s.repository.InsertMany(ctx, entities)
	}

}

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
