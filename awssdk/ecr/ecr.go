package ecr

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *ecr.Client
)

func init() {
	svc = ecr.NewFromConfig(awsenv.Cfg)

}

func DescribeRepositories() {
	input := &ecr.DescribeRepositoriesInput{}
	output, err := svc.DescribeRepositories(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, repository := range output.Repositories {
		fmt.Println("RepositoryName:", *repository.RepositoryName, "\tRepositoryUri:", *repository.RepositoryUri)
	}
}
