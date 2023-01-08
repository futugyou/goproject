package ecr

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/tools"
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
			return output.ImageDetails[i].ImagePushedAt.After(*output.ImageDetails[j].ImagePushedAt)
		})

		for i := 0; i < len(output.ImageDetails) && i < 5; i++ {
			fmt.Println("\tImageDigest:", *output.ImageDetails[i].ImageDigest, output.ImageDetails[i].ImagePushedAt)
		}

		imageIdentifier := types.ImageIdentifier{ImageDigest: output.ImageDetails[0].ImageDigest}
		batchinput := &ecr.BatchGetImageInput{
			ImageIds:       []types.ImageIdentifier{imageIdentifier},
			RepositoryName: repository.RepositoryName,
		}

		batchoutput, err := svc.BatchGetImage(awsenv.EmptyContext, batchinput)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, image := range batchoutput.Images {
			tag := tools.Sha1(*image.ImageId.ImageTag)
			input := &ecr.PutImageInput{
				ImageManifest:  image.ImageManifest,
				RepositoryName: image.RepositoryName,
				ImageTag:       &tag,
			}

			output, err := svc.PutImage(awsenv.EmptyContext, input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("\t\tImageTag:", *output.Image.ImageId.ImageTag)
			break
		}

		fmt.Println()
	}
}

func CreateRepository() {
	input := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String("jenkins"),
	}

	output, err := svc.CreateRepository(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("RepositoryUri:", *output.Repository.RepositoryUri)
}

func DeleteRepository() {
	input := &ecr.DeleteRepositoryInput{
		RepositoryName: aws.String("jenkins"),
	}

	output, err := svc.DeleteRepository(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("RepositoryUri:", *output.Repository.RepositoryUri)
}
