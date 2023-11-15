package s3

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *s3.Client
)

func init() {
	svc = s3.NewFromConfig(awsenv.Cfg)
}

func ListBuckets() {
	input := s3.ListBucketsInput{}
	output, err := svc.ListBuckets(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, bucket := range output.Buckets {
		fmt.Println(*bucket.Name, bucket.CreationDate)
	}
}

func GetBucketCors(bucketName string) {
	input := s3.GetBucketCorsInput{
		Bucket: &bucketName,
	}
	output, err := svc.GetBucketCors(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, rule := range output.CORSRules {
		fmt.Println(rule.AllowedHeaders, rule.AllowedMethods, rule.AllowedOrigins, rule.ExposeHeaders, rule.MaxAgeSeconds)
	}
}

func GetBucketPolicy(bucketName string) {
	input := s3.GetBucketPolicyInput{
		Bucket: &bucketName,
	}
	output, err := svc.GetBucketPolicy(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if output.Policy != nil {
		fmt.Println(*output.Policy)
	}
}

func GetBucketPolicyStatus(bucketName string) {
	input := s3.GetBucketPolicyStatusInput{
		Bucket: &bucketName,
	}
	output, err := svc.GetBucketPolicyStatus(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if output.PolicyStatus != nil {
		fmt.Println(output.PolicyStatus.IsPublic)
	}
}

func ListObjectsV2(bucketName string) {
	objInput := s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}
	ojbOutput, err := svc.ListObjectsV2(awsenv.EmptyContext, &objInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("name:", *ojbOutput.Name)
	for _, obj := range ojbOutput.Contents {
		fmt.Println("\tKey:", *obj.Key, "\tLastModified:", obj.LastModified)
	}
	for _, obj := range ojbOutput.CommonPrefixes {
		fmt.Println("\tPrefix:", *obj.Prefix)
	}
}

func GetObject(bucket, key string) {
	input := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	output, err := svc.GetObject(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	outFile, err := os.Create("./test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	// handle err
	defer outFile.Close()
	_, err = io.Copy(outFile, output.Body)
	if err != nil {
		fmt.Println(err)
	}
}

func PutObject(bucket, key string) {
	outFile, err := os.Open("./test.tx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()

	input := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   outFile,
	}
	output, err := svc.PutObject(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}
