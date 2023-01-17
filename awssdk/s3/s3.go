package s3

import (
	"fmt"

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
		fmt.Println(*bucket.Name)
	}
}
