package efs

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *efs.Client
)

func init() {
	svc = efs.NewFromConfig(awsenv.Cfg)
}

func DescribeFileSystems() {
	input := efs.DescribeFileSystemsInput{}
	output, err := svc.DescribeFileSystems(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, fs := range output.FileSystems {
		fmt.Println(
			// *fs.AvailabilityZoneId,
			// *fs.AvailabilityZoneName,
			*fs.CreationTime,
			*fs.CreationToken,
			*fs.Encrypted,
			*fs.FileSystemArn,
			*fs.FileSystemId,
			*fs.KmsKeyId,
			fs.LifeCycleState,
			*fs.Name,
			fs.NumberOfMountTargets,
			*fs.OwnerId,
			fs.PerformanceMode,
			// *fs.ProvisionedThroughputInMibps,
			fs.ThroughputMode,
		)
	}
}
