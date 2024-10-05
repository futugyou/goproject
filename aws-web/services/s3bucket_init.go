package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

func (s *S3bucketService) InitData(ctx context.Context) {
	log.Println("s3 sync start..")

	accountService := NewAccountService()
	accounts := accountService.GetAllAccounts(ctx)

	entities := make([]entity.S3bucketEntity, 0)
	items := make([]entity.S3bucketItemEntity, 0)
	for _, account := range accounts {
		if !account.Valid {
			continue
		}
		awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
		buckets, err := s.GetAllS3Bucket(ctx)
		if err != nil {
			continue
		}

		for _, bucket := range buckets {
			b := entity.S3bucketEntity{
				Id:           account.Id + *bucket.Name,
				AccountId:    account.Id,
				Name:         *bucket.Name,
				Region:       s.GetBucketRegion(ctx, bucket.Name, account.Region),
				IsPublic:     s.GetBucketPolicyStatus(ctx, bucket.Name),
				Policy:       s.GetBucketPolicy(ctx, bucket.Name),
				Permissions:  s.GetBucketPermissions(ctx, bucket.Name),
				CreationDate: *bucket.CreationDate,
			}
			entities = append(entities, b)

			objs, err := s.ListItemsByBucketName(ctx, bucket.Name)
			if err != nil {
				continue
			}

			for _, obj := range objs {
				item := entity.S3bucketItemEntity{
					Id:           account.Id + *bucket.Name + *obj.Key,
					BucketName:   *bucket.Name,
					Key:          *obj.Key,
					Size:         *obj.Size,
					CreationDate: *obj.LastModified,
				}
				items = append(items, item)
			}
		}
	}

	if len(entities) > 0 {
		s.repository.DeleteAll(ctx)
		s.repository.InsertMany(ctx, entities)
	}

	if len(items) > 0 && len(entities) > 0 {
		for _, b := range entities {
			sunItem := FilterS3bucketItemEntity(items, b.Name)
			if len(sunItem) == 0 {
				continue
			}

			s.itemRepository.DeleteByBucketName(ctx, b.Name)
			s.itemRepository.InsertMany(ctx, sunItem)
		}
	}

	log.Println("total bucket count ", len(entities))
	log.Println("total item count ", len(items))
	log.Println("s3 sync end.")
}

func FilterS3bucketItemEntity(items []entity.S3bucketItemEntity, name string) []entity.S3bucketItemEntity {
	sub := make([]entity.S3bucketItemEntity, 0)
	for _, i := range items {
		if i.BucketName == name {
			sub = append(sub, i)
		}
	}
	return sub
}

func (s *S3bucketService) GetBucketRegion(ctx context.Context, name *string, region string) string {
	svc := s3.NewFromConfig(awsenv.Cfg)

	input := s3.GetBucketLocationInput{
		Bucket: name,
	}
	output, err := svc.GetBucketLocation(ctx, &input)
	if err != nil {
		log.Println("GetBucketLocation error.")
		return region
	}

	return string(output.LocationConstraint)
}

func (s *S3bucketService) GetBucketPolicyStatus(ctx context.Context, name *string) bool {
	svc := s3.NewFromConfig(awsenv.Cfg)

	input := s3.GetBucketPolicyStatusInput{
		Bucket: name,
	}
	output, err := svc.GetBucketPolicyStatus(ctx, &input)
	if err != nil {
		log.Println("GetBucketPolicyStatus error.")
		return false
	}

	if output.PolicyStatus != nil && output.PolicyStatus.IsPublic != nil {
		return *output.PolicyStatus.IsPublic
	}

	return false
}

func (s *S3bucketService) GetBucketPolicy(ctx context.Context, name *string) string {
	svc := s3.NewFromConfig(awsenv.Cfg)

	input := s3.GetBucketPolicyInput{
		Bucket: name,
	}
	output, err := svc.GetBucketPolicy(ctx, &input)
	if err != nil {
		log.Println("GetBucketPolicy error.")
		return ""
	}

	if output.Policy != nil {
		return *output.Policy
	}

	return ""
}

func (s *S3bucketService) GetBucketPermissions(ctx context.Context, name *string) []string {
	svc := s3.NewFromConfig(awsenv.Cfg)
	input := s3.GetBucketAclInput{
		Bucket: name,
	}
	output, err := svc.GetBucketAcl(ctx, &input)
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

func (s *S3bucketService) GetAllS3Bucket(ctx context.Context) ([]types.Bucket, error) {
	svc := s3.NewFromConfig(awsenv.Cfg)
	input := s3.ListBucketsInput{}

	output, err := svc.ListBuckets(ctx, &input)
	if err != nil {
		log.Println("get ListBuckets error.")
		return nil, err
	}

	return output.Buckets, nil
}

func (s *S3bucketService) ListItemsByBucketName(ctx context.Context, name *string) ([]types.Object, error) {
	svc := s3.NewFromConfig(awsenv.Cfg)

	objInput := s3.ListObjectsV2Input{
		Bucket: name,
	}

	output, err := svc.ListObjectsV2(ctx, &objInput)
	if err != nil {
		log.Println("get ListObjectsV2 error.")
		return nil, err
	}

	return output.Contents, nil
}

func (s *S3bucketService) ListItems(ctx context.Context, name string, perfix string, del string) (*s3.ListObjectsV2Output, error) {
	svc := s3.NewFromConfig(awsenv.Cfg)

	objInput := s3.ListObjectsV2Input{
		Bucket: aws.String(name),
		// Delimiter: aws.String("/"),
	}

	if del != "1" {
		objInput.Delimiter = aws.String("/")
	}

	if len(perfix) > 0 {
		objInput.Prefix = aws.String(perfix)
	}

	output, err := svc.ListObjectsV2(ctx, &objInput)
	if err != nil {
		log.Println("get ListObjectsV2 error.")
		return nil, err
	}

	return output, nil
}

func (s *S3bucketService) GetS3Object(ctx context.Context, bucket string, key string) (*s3.GetObjectOutput, error) {
	svc := s3.NewFromConfig(awsenv.Cfg)

	input := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	return svc.GetObject(ctx, &input)
}

func (s *S3bucketService) PresignGetObject(ctx context.Context, bucket string, key string) string {
	svc := s3.NewFromConfig(awsenv.Cfg)
	presignClient := s3.NewPresignClient(svc)
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(600 * int64(time.Second))
	})
	if err != nil {
		fmt.Println("PresignGetObject error")
		return ""
	}
	return request.URL
}
