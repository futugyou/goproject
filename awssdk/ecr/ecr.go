package ecr

import (
	"fmt"
	"sort"

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

		// DescribeImages
		input := &ecr.DescribeImagesInput{
			RepositoryName: repository.RepositoryName,
		}
		output, err := svc.DescribeImages(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(output.ImageDetails) == 0 {
			continue
		}

		sort.Slice(output.ImageDetails, func(i, j int) bool {
			return output.ImageDetails[i].ImagePushedAt.Before(*output.ImageDetails[j].ImagePushedAt)
		})

		for i := 0; i < len(output.ImageDetails); i++ {
			fmt.Println("\tImageDigest:", *output.ImageDetails[i].ImageDigest, output.ImageDetails[i].ImagePushedAt)
		}
		fmt.Println()
	}
}
