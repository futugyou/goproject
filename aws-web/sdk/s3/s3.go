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

func ListBuckets() []string {
	input := s3.ListBucketsInput{}
	result := make([]string, 0)
	output, err := svc.ListBuckets(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return result
	}

	for _, bucket := range output.Buckets {
		fmt.Println(*bucket.Name, bucket.CreationDate)
		result = append(result, *bucket.Name)
	}

	return result
}

// error if not exist
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

// error if not exist
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
		fmt.Println(bucketName, "\t", *output.Policy)
	}
}

// error if not exist
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
		fmt.Println(bucketName, "\t", output.PolicyStatus.IsPublic)
	}
}

// error if not unsupported
// no need
func GetBucketAccelerateConfiguration(bucketName string) {
	input := s3.GetBucketAccelerateConfigurationInput{
		Bucket: &bucketName,
	}
	output, err := svc.GetBucketAccelerateConfiguration(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bucketName, "\t", output.Status)
}

func GetBucketAcl(bucketName string) {
	input := s3.GetBucketAclInput{
		Bucket: &bucketName,
	}
	output, err := svc.GetBucketAcl(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if output.Owner != nil {
		if output.Owner.DisplayName != nil {
			fmt.Println(*output.Owner.DisplayName)
		}
		if output.Owner.ID != nil {
			fmt.Println(*output.Owner.ID)
		}
	}
	for _, v := range output.Grants {
		fmt.Println(v.Permission)
	}
}

func GetBucketLocation(bucketName string) {
	input := s3.GetBucketLocationInput{
		Bucket: &bucketName,
	}
	output, err := svc.GetBucketLocation(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bucketName, "\t", output.LocationConstraint)
}

// no need
func GetBucketWebsite(bucketName string) {
	input := s3.GetBucketWebsiteInput{
		Bucket: &bucketName,
	}
	output, err := svc.GetBucketWebsite(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bucketName, "\t", output.IndexDocument)
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
