package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
)

type RegionService struct {
}

func NewRegionService() *RegionService {
	return &RegionService{}
}

func (s *RegionService) GetRegions(ctx context.Context) ([]model.AwsRegion, error) {
	input := &ec2.DescribeRegionsInput{}
	svc := ec2.NewFromConfig(awsenv.Cfg)
	output, err := svc.DescribeRegions(ctx, input)
	if err != nil {
		return nil, err
	}

	result := make([]model.AwsRegion, 0)

	for _, r := range output.Regions {
		region := model.AwsRegion{
			Endpoint: *r.Endpoint,
			Status:   *r.OptInStatus,
			Region:   *r.RegionName,
		}

		result = append(result, region)
	}

	return result, nil
}
