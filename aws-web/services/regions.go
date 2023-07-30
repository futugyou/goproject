package services

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

type AwsRegion struct {
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
	Region   string `json:"region"`
}

type RegionService struct {
}

func NewRegionService() *RegionService {
	return &RegionService{}
}

func (s *RegionService) GetAllRegionInCurrentAccount() ([]AwsRegion, error) {
	input := &ec2.DescribeRegionsInput{}
	// TODO: aws config
	svc := ec2.NewFromConfig(awsenv.Cfg)
	output, err := svc.DescribeRegions(awsenv.EmptyContext, input)
	if err != nil {
		return nil, err
	}

	result := make([]AwsRegion, 0)

	for _, r := range output.Regions {
		region := AwsRegion{
			Endpoint: *r.Endpoint,
			Status:   *r.OptInStatus,
			Region:   *r.RegionName,
		}

		result = append(result, region)
	}

	return result, nil
}
