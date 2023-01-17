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
