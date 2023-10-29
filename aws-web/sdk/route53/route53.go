package route53

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *route53.Client
)

func init() {
	svc = route53.NewFromConfig(awsenv.Cfg)
}

func GetHostedZone() {
	input := &route53.GetHostedZoneInput{
		Id: aws.String(awsenv.HostedZoneId),
	}

	result, err := svc.GetHostedZone(awsenv.EmptyContext, input)
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	for _, vpc := range result.VPCs {
		fmt.Println(*vpc.VPCId, vpc.VPCRegion.Values())
	}
}

func GetHostedZoneVpcId(hostedZoneId string) string {
	svc = route53.NewFromConfig(awsenv.Cfg)
	input := &route53.GetHostedZoneInput{
		Id: aws.String(hostedZoneId),
	}

	result, err := svc.GetHostedZone(awsenv.EmptyContext, input)
	if err != nil {
		return ""
	}

	for _, vpc := range result.VPCs {
		if vpc.VPCId != nil {
			return *vpc.VPCId
		}
	}

	return ""
}
