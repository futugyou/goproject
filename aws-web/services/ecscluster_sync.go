package services

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

func (e *EcsClusterService) SyncAllEcsServices(ctx context.Context) {
	log.Println("start..")
	accountService := NewAccountService()
	accounts := accountService.GetAllAccounts(ctx)
	services := make([]entity.EcsServiceEntity, 0)
	for _, account := range accounts {
		awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
		svc := ecs.NewFromConfig(awsenv.Cfg)
		clusterinput := &ecs.ListClustersInput{}
		clusteroutput, err := svc.ListClusters(ctx, clusterinput)
		if err != nil {
			log.Println("list ecs clusters error")
			continue
		}

		for _, cluster := range clusteroutput.ClusterArns {
			serviceinput := &ecs.ListServicesInput{
				Cluster:    aws.String(cluster),
				MaxResults: aws.Int32(100),
			}

			serviceoutput, err := svc.ListServices(ctx, serviceinput)
			if err != nil {
				log.Println("list ecs service error")
				continue
			}

			var serviceArns []string = serviceoutput.ServiceArns
			for {
				if len(serviceArns) == 0 {
					break
				}

				t := serviceArns
				if len(t) > 10 {
					t = serviceArns[:10]
				}

				describeinput := &ecs.DescribeServicesInput{
					Cluster:  aws.String(cluster),
					Services: t,
				}

				describeoutput, err := svc.DescribeServices(ctx, describeinput)
				if len(serviceArns) > 10 {
					serviceArns = serviceArns[10:]
				} else {
					serviceArns = []string{}
				}

				if err != nil {
					log.Println("dscribe ecs serive error")
					continue
				}

				for _, v := range describeoutput.Services {
					t := time.Now()
					if len(v.Deployments) > 0 && v.Deployments[0].UpdatedAt != nil {
						t = *v.Deployments[0].UpdatedAt
					}

					tmp := strings.Split(cluster, "/")
					entity := entity.EcsServiceEntity{
						AccountId:      account.Id,
						Cluster:        tmp[len(tmp)-1],
						ClusterArn:     *v.ClusterArn,
						ServiceName:    *v.ServiceName,
						ServiceNameArn: *v.ServiceArn,
						RoleArn:        *v.RoleArn,
						OperateAt:      t,
					}

					services = append(services, entity)
				}
			}
		}
	}

	log.Println("get finish, count: ", len(services))
	err := e.repository.BulkWrite(ctx, services)
	log.Println("ecs write finish: ", err)
}
