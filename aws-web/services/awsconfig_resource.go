package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
	c "github.com/futugyousuzu/goproject/awsgolang/viewmodel/awsconfigConfiguration"

	"github.com/futugyousuzu/goproject/awsgolang/sdk/ecs"
	"github.com/futugyousuzu/goproject/awsgolang/sdk/iam"
	"github.com/futugyousuzu/goproject/awsgolang/sdk/loadbalancing"
	"github.com/futugyousuzu/goproject/awsgolang/sdk/route53"
	"github.com/futugyousuzu/goproject/awsgolang/sdk/servicediscovery"
	"golang.org/x/exp/slices"
)

func CreateAwsConfigEntity(data model.AwsConfigRawData, vpcinfos []c.VpcInfo) entity.AwsConfigEntity {
	configuration := getDataString(data.Configuration)
	name := getName(data.ResourceID, data.ResourceName, data.Tags)
	vpcid, subnetId, subnetIds, securityGroups, availabilityZone := getVpcInfo(data.ResourceType, configuration, vpcinfos)
	loginURL, loggedInURL := createConsoleUrls(data)
	resource := entity.AwsConfigEntity{
		ID:                           getId(data.ARN, data.ResourceID),
		Label:                        name,
		AccountID:                    data.AwsAccountID,
		Arn:                          data.ARN,
		AvailabilityZone:             data.AvailabilityZone,
		AwsRegion:                    data.AwsRegion,
		Configuration:                configuration,
		ConfigurationItemCaptureTime: data.ConfigurationItemCaptureTime,
		ConfigurationItemStatus:      data.ConfigurationItemStatus,
		ConfigurationStateID:         data.ConfigurationStateID,
		ResourceCreationTime:         data.ResourceCreationTime,
		ResourceID:                   data.ResourceID,
		ResourceName:                 data.ResourceName,
		ResourceType:                 data.ResourceType,
		Tags:                         getDataString(data.Tags),
		Version:                      data.ConfigurationItemVersion,
		VpcID:                        vpcid,
		SubnetID:                     subnetId,
		SubnetIds:                    subnetIds,
		Title:                        name,
		SecurityGroups:               securityGroups,
		LoginURL:                     loginURL,
		LoggedInURL:                  loggedInURL,
	}

	if len(availabilityZone) > 0 {
		resource.AvailabilityZone = availabilityZone
	}

	return resource
}

func AddIndividualResource(resources []entity.AwsConfigEntity, vpcinfos []c.VpcInfo) []entity.AwsConfigEntity {
	discoverys := make([]entity.AwsConfigEntity, 0)
	ecscluster := make([]entity.AwsConfigEntity, 0)

	for _, resource := range resources {
		switch resource.ResourceType {
		case "AWS::ServiceDiscovery::Service":
			discoverys = append(discoverys, resource)
		case "AWS::ECS::Cluster":
			ecscluster = append(ecscluster, resource)
		}
	}

	// 1. Service Discovery Namespace
	namespaces := createServiceDiscoveryNamespaces(discoverys)
	resources = append(resources, namespaces...)
	// 2. ecs task
	ecsTasks := createEcsTaskResources(ecscluster, vpcinfos)
	resources = append(resources, ecsTasks...)
	// 3. aws managed polices
	polices := createAwsManagedPolices()
	resources = append(resources, polices...)
	// 4. target groups
	targetGroups := createTargetGroups()
	resources = append(resources, targetGroups...)

	return resources
}

func createTargetGroups() []entity.AwsConfigEntity {
	result := make([]entity.AwsConfigEntity, 0)
	tgs := loadbalancing.GetTargetGroups()
	for _, tg := range tgs {
		confString, _ := json.Marshal(tg)
		_, _, _, accid := ParseARN(*tg.TargetGroupArn)
		result = append(result, entity.AwsConfigEntity{
			ID:                           *tg.TargetGroupArn,
			Label:                        *tg.TargetGroupName,
			AccountID:                    accid,
			Arn:                          *tg.TargetGroupArn,
			AvailabilityZone:             "Multiple Availability Zones",
			AwsRegion:                    awsenv.Cfg.Region,
			Configuration:                string(confString),
			ConfigurationItemCaptureTime: time.Time{},
			ConfigurationItemStatus:      "",
			ConfigurationStateID:         "0",
			ResourceCreationTime:         time.Time{},
			ResourceID:                   *tg.TargetGroupArn,
			ResourceName:                 *tg.TargetGroupName,
			ResourceType:                 "AWS::ElasticLoadBalancingV2::TargetGroup",
			Tags:                         "",
			Version:                      "",
			VpcID:                        *tg.VpcId,
			SubnetID:                     "",
			SubnetIds:                    []string{},
			Title:                        *tg.TargetGroupName,
			SecurityGroups:               []string{},
			LoginURL:                     "",
			LoggedInURL:                  "",
		})
	}
	return result
}

func createAwsManagedPolices() []entity.AwsConfigEntity {
	result := make([]entity.AwsConfigEntity, 0)
	polices := iam.ListAttachedAwsManagedPolices()
	for _, policy := range polices {
		tags := ""
		if len(policy.Tags) > 0 {
			tags = getDataString(policy.Tags)
		}
		confString, _ := json.Marshal(policy)
		result = append(result, entity.AwsConfigEntity{
			ID:                           *policy.Arn,
			Label:                        *policy.PolicyName,
			AccountID:                    "aws",
			Arn:                          *policy.Arn,
			AvailabilityZone:             "Not Applicable",
			AwsRegion:                    awsenv.Cfg.Region,
			Configuration:                string(confString),
			ConfigurationItemCaptureTime: *policy.CreateDate,
			ConfigurationItemStatus:      "",
			ConfigurationStateID:         "0",
			ResourceCreationTime:         *policy.CreateDate,
			ResourceID:                   *policy.Arn,
			ResourceName:                 *policy.PolicyName,
			ResourceType:                 "AWS::IAM::AWSManagedPolicy",
			Tags:                         tags,
			Version:                      "",
			VpcID:                        "",
			SubnetID:                     "",
			SubnetIds:                    []string{},
			Title:                        *policy.PolicyName,
			SecurityGroups:               []string{},
			LoginURL:                     "",
			LoggedInURL:                  "",
		})
	}
	return result
}

func createEcsTaskResources(ecscluster []entity.AwsConfigEntity, vpcinfos []c.VpcInfo) []entity.AwsConfigEntity {
	result := make([]entity.AwsConfigEntity, 0)
	clusterNames := make([]string, 0)
	for _, cluster := range ecscluster {
		clusterNames = append(clusterNames, cluster.ResourceName)
	}

	tasks := ecs.GetEcsTasksByCluster(clusterNames)
	for _, task := range tasks {
		i := slices.IndexFunc(ecscluster, func(cluster entity.AwsConfigEntity) bool {
			return cluster.Arn == *task.ClusterArn
		})
		if i == -1 {
			continue
		}
		cluster := ecscluster[i]
		confString, _ := json.Marshal(task)
		vpcid, subnetId, subnetIds, securityGroups, availabilityZone := getVpcInfo("AWS::ECS::TASK", string(confString), vpcinfos)
		tags := ""
		if len(task.Tags) > 0 {
			tags = getDataString(task.Tags)
		}
		taskEntity := entity.AwsConfigEntity{
			ID:                           *task.TaskArn,
			Label:                        *task.Group,
			AccountID:                    cluster.AccountID,
			Arn:                          *task.TaskArn,
			AvailabilityZone:             availabilityZone,
			AwsRegion:                    cluster.AwsRegion,
			Configuration:                string(confString),
			ConfigurationItemCaptureTime: *task.CreatedAt,
			ConfigurationItemStatus:      string(task.HealthStatus),
			ConfigurationStateID:         "0",
			ResourceCreationTime:         *task.CreatedAt,
			ResourceID:                   *task.TaskArn,
			ResourceName:                 *task.Group,
			ResourceType:                 "AWS::ECS::TASK",
			Tags:                         tags,
			Version:                      "",
			VpcID:                        vpcid,
			SubnetID:                     subnetId,
			SubnetIds:                    subnetIds,
			Title:                        *task.Group,
			SecurityGroups:               securityGroups,
			LoginURL:                     "",
			LoggedInURL:                  "",
		}
		result = append(result, taskEntity)
	}
	return result
}

func createServiceDiscoveryNamespaces(discoverys []entity.AwsConfigEntity) []entity.AwsConfigEntity {
	namespaceEntityList := make([]entity.AwsConfigEntity, 0)
	namespaceList := make([]string, 0)
	for _, discovery := range discoverys {
		var configuration c.ServiceDiscoveryConfiguration
		err := json.Unmarshal([]byte(discovery.Configuration), &configuration)
		if err != nil {
			log.Println(err)
			continue
		}

		namespaceid := configuration.NamespaceID
		if slices.Contains(namespaceList, namespaceid) {
			continue
		}

		namespaceList = append(namespaceList, namespaceid)
		namespace := servicediscovery.GetNamespaceDetail(namespaceid)

		if namespace == nil {
			continue
		}

		namespaceEntity := entity.AwsConfigEntity{
			ID:                           *namespace.Arn,
			Label:                        *namespace.Name,
			AccountID:                    discovery.AccountID,
			Arn:                          *namespace.Arn,
			AvailabilityZone:             discovery.AvailabilityZone,
			AwsRegion:                    discovery.AwsRegion,
			Configuration:                "{}",
			ConfigurationItemCaptureTime: *namespace.CreateDate,
			ConfigurationItemStatus:      "",
			ConfigurationStateID:         "0",
			ResourceCreationTime:         *namespace.CreateDate,
			ResourceID:                   *namespace.Id,
			ResourceName:                 *namespace.Name,
			ResourceType:                 "AWS::ServiceDiscovery::Namespace",
			Tags:                         "",
			Version:                      "",
			VpcID:                        "",
			SubnetID:                     "",
			SubnetIds:                    []string{},
			Title:                        *namespace.Name,
			SecurityGroups:               []string{},
			LoginURL:                     "",
			LoggedInURL:                  "",
		}

		if namespace.Properties != nil &&
			namespace.Properties.DnsProperties != nil &&
			namespace.Properties.DnsProperties.HostedZoneId != nil {
			vpcid := route53.GetHostedZoneVpcId(*namespace.Properties.DnsProperties.HostedZoneId)
			namespaceEntity.VpcID = vpcid
			discovery.VpcID = vpcid
		}

		namespaceEntityList = append(namespaceEntityList, namespaceEntity)

	}
	return namespaceEntityList
}

func GetAllVpcInfos(datas []model.AwsConfigRawData) []c.VpcInfo {
	vpcInfos := make([]c.VpcInfo, 0)
	for _, data := range datas {
		if data.ResourceType == "AWS::EC2::Subnet" {
			var config c.SubnetConfiguration
			json.Unmarshal([]byte(getDataString(data.Configuration)), &config)

			vpcid := config.VpcID
			subnetId := config.SubnetID
			availabilityZone := config.AvailabilityZone

			found := false
			for i := 0; i < len(vpcInfos); i++ {
				if vpcid == vpcInfos[i].VpcId {
					found = true
					vpcInfos[i].Subnets = append(vpcInfos[i].Subnets, c.SubnetInfo{
						Subnet:           subnetId,
						AvailabilityZone: availabilityZone,
					})
				}
			}

			if !found {
				vpc := c.VpcInfo{
					VpcId: vpcid,
					Subnets: []c.SubnetInfo{{
						Subnet:           subnetId,
						AvailabilityZone: availabilityZone,
					}},
				}
				vpcInfos = append(vpcInfos, vpc)
			}
		}

	}
	return vpcInfos
}

func ParseARN(arn string) (partition, service, region, accountId string) {
	partition = ""
	service = ""
	region = ""
	accountId = ""
	segments := strings.Split(arn, ":")
	if len(segments) > 4 {
		partition = segments[1]
		service = segments[2]
		region = segments[3]
		accountId = segments[4]
	}
	return
}

func getId(arn string, resourceID string) string {
	if len(arn) == 0 {
		return resourceID
	}
	return arn
}

func getName(resourceID string, resourceName string, tags []model.Tag) string {
	if len(tags) > 0 {
		for _, tag := range tags {
			if tag.Key == "Name" && len(tag.Value) > 0 {
				return tag.Value
			}
		}
	}

	if len(resourceName) != 0 {
		return resourceName
	}

	return resourceID
}

func getDataString(con interface{}) string {
	if con == nil {
		return "{}"
	} else {
		d, _ := json.Marshal(con)
		return string(d)
	}
}

func getVpcInfo(resourceType string, configuration string, vpcinfos []c.VpcInfo) (vpcid string, subnetId string, subnetIds []string, securityGroups []string, availabilityZone string) {
	vpcid = ""
	subnetId = ""
	subnetIds = make([]string, 0)
	securityGroups = make([]string, 0)
	availabilityZone = ""

	foundVpn := func(ids []string) string {
		for _, vpn := range vpcinfos {
			for _, net := range vpn.Subnets {
				if slices.Contains(ids, net.Subnet) {
					return vpn.VpcId
				}
			}
		}

		return ""
	}

	foundAvailabilityZone := func(ids []string) string {
		azs := make([]string, 0)
		for _, vpn := range vpcinfos {
			for _, net := range vpn.Subnets {
				if !slices.Contains(ids, net.Subnet) ||
					slices.Contains(azs, net.AvailabilityZone) {
					continue
				}

				azs = append(azs, net.AvailabilityZone)
			}
		}

		if len(azs) > 0 {
			slices.Sort(azs)
			return strings.Join(azs, ",")
		}

		return ""
	}
	switch resourceType {
	case "AWS::EC2::VPCEndpoint":
		var config c.VPCEndpointConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
			subnetIds = config.SubnetIDS
			for _, v := range config.Groups {
				securityGroups = append(securityGroups, v.GroupID)
			}
		}
	case "AWS::EC2::VPC":
		var config c.VPCConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
		}
	case "AWS::EC2::Subnet":
		var config c.SubnetConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
			subnetId = config.SubnetID
		}
	case "AWS::AmazonMQ::Broker":
		var config c.AmazonMQConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			subnetIds = config.SubnetIDS
			securityGroups = config.SecurityGroups
		}
	case "AWS::EC2::NatGateway":
		var config c.NatGatewayConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			subnetId = config.SubnetID
			vpcid = config.VpcID
		}
	case "AWS::RDS::DBInstance":
		var config c.DBInstanceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			for _, v := range config.VpcSecurityGroups {
				securityGroups = append(securityGroups, v.VpcSecurityGroupID)
			}
			vpcid = config.DBSubnetGroup.VpcID
			for _, v := range config.DBSubnetGroup.Subnets {
				subnetIds = append(subnetIds, v.SubnetIdentifier)
			}
		}
	case "AWS::EC2::SecurityGroup":
		var config c.SecurityGroupConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
			securityGroups = []string{config.GroupID}
		}
	case "AWS::EC2::NetworkInterface":
		var config c.NetworkInterfaceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
			subnetId = config.SubnetID
			for _, v := range config.Groups {
				securityGroups = append(securityGroups, v.GroupID)
			}
		}
	case "AWS::Redshift::ClusterSubnetGroup":
		var config c.RedshiftConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID

			for _, v := range config.Subnets {
				subnetIds = append(subnetIds, v.SubnetIdentifier)
			}
		}
	case "AWS::ElasticLoadBalancingV2::LoadBalancer":
		var config c.LoadBalancerConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
			for _, v := range config.SecurityGroups {
				securityGroups = append(securityGroups, v.Value)
			}

			for _, v := range config.AvailabilityZones {
				subnetIds = append(subnetIds, v.SubnetID)
			}
		}
	case "AWS::ECS::Service":
		var config c.ECSServiceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			subnetIds = config.NetworkConfiguration.AwsvpcConfiguration.Subnets
			securityGroups = config.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups
		}
	case "AWS::EC2::NetworkAcl":
		var config c.NetworkAclConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
			for _, v := range config.Associations {
				subnetIds = append(subnetIds, v.SubnetID)
			}
		}
	case "AWS::Lambda::Function":
		var config c.FunctionConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			subnetIds = config.VpcConfig.SubnetIDS
			securityGroups = config.VpcConfig.SecurityGroupIDS
		}
	case "AWS::EC2::RouteTable":
		var config c.RouteTableConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
		}
	case "AWS::EC2::Instance":
		var config c.Ec2Configuration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			vpcid = config.VpcID
			subnetId = config.SubnetID
			for _, sg := range config.SecurityGroups {
				securityGroups = append(securityGroups, sg.GroupID)
			}
		}
	case "AWS::ECS::TASK":
		var config c.ECSTaskConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err != nil {
			log.Println(err)
		} else {
			for _, att := range config.Attachments {
				for _, d := range att.Details {
					if d.Name == "subnetId" && len(d.Value) > 0 {
						subnetIds = append(subnetIds, d.Value)
					}
				}
			}
		}
	}

	slices.Sort(subnetIds)

	if vpcid == "" {
		if len(subnetIds) > 0 {
			vpcid = foundVpn(subnetIds)
		} else if len(subnetId) > 0 {
			vpcid = foundVpn([]string{subnetId})
		}
	}

	if len(subnetId) == 0 && len(subnetIds) == 1 {
		subnetId = subnetIds[0]
	}

	if len(subnetIds) > 1 {
		availabilityZone = foundAvailabilityZone(subnetIds)
		subnetId = ""
	}

	return
}

func createSignInHostname(accountId string, service string) string {
	return fmt.Sprintf("https://%s.signin.aws.amazon.com/console/%s", accountId, service)
}

func createLoggedInHostname(awsRegion string, service string) string {
	return fmt.Sprintf(`https://%s.console.aws.amazon.com%s/home`, awsRegion, service)
}

func createConsoleUrls(resource model.AwsConfigRawData) (loginURL string, loggedInURL string) {
	resourceType := resource.ResourceType
	resourceName := resource.ResourceName
	accountId := resource.AwsAccountID
	awsRegion := resource.AwsRegion
	loginURL = ""
	loggedInURL = ""
	switch resourceType {
	case "AWS::Lambda::Function":
		loginURL =
			fmt.Sprintf(`%s?region=%s#/functions/%s?tab=graph`,
				createSignInHostname(accountId, "lambda"),
				awsRegion,
				resourceName,
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#/functions/%s?tab=graph`,
				createLoggedInHostname(awsRegion, "lambda"),
				awsRegion,
				resourceName,
			)
	case "AWS::IAM::Policy":
		loginURL =
			fmt.Sprintf(`%s?home?#%s`,
				createSignInHostname(accountId, "iam"),
				"/policies",
			)
		loggedInURL =
			fmt.Sprintf(`https://console.aws.amazon.com/%s/home?#%s`,
				"iam",
				"/policies",
			)
	case "AWS::S3::Bucket":
		loginURL =
			fmt.Sprintf(`%s?bucket=%s`,
				createSignInHostname(accountId, "s3"),
				resourceName,
			)
		loggedInURL =
			fmt.Sprintf(`https://s3.console.aws.amazon.com/s3/buckets/%s/?region=%s`,
				resourceName,
				awsRegion,
			)
	case "AWS::EC2::VPC":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"vpcs:sort=VpcId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"vpcs:sort=VpcId",
			)
	case "AWS::EC2::NetworkInterface":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"NIC:sort=description",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"NIC:sort=description",
			)
	case "AWS::EC2::Instance":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"Instances:sort=instanceId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"Instances:sort=instanceId",
			)
	case "AWS::EC2::Volume":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"Volumes:sort=desc:name",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"Volumes:sort=desc:name",
			)
	case "AWS::EC2::Subnet":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"subnets:sort=SubnetId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"subnets:sort=SubnetId",
			)
	case "AWS::EC2::SecurityGroup":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"SecurityGroups:sort=groupId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"SecurityGroups:sort=groupId",
			)
	case "AWS::EC2::RouteTable":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"RouteTables:sort=routeTableId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"RouteTables:sort=routeTableId",
			)
	case "AWS::EC2::InternetGateway":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"igws:sort=internetGatewayId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"igws:sort=internetGatewayId",
			)
	case "AWS::EC2::NetworkAcl":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"acls:sort=networkAclId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"acls:sort=networkAclId",
			)
	case "AWS::ElasticLoadBalancingV2::LoadBalancer":
	case "AWS::ElasticLoadBalancingV2::Listener":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"LoadBalancers:",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"LoadBalancers:",
			)

	case "AWS::ElasticLoadBalancingV2::TargetGroup":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"TargetGroups:",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"TargetGroups:",
			)
	case "AWS::EC2::EIP":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"Addresses:sort=PublicIp",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"Addresses:sort=PublicIp",
			)
	}

	return
}

func FilterResource(datas []model.AwsConfigRawData) []model.AwsConfigRawData {
	resuls := make([]model.AwsConfigRawData, 0)
	for _, d := range datas {
		if d.ResourceType == "AWS::ACM::Certificate" ||
			d.ResourceType == "AWS::AmazonMQ::Broker" ||
			d.ResourceType == "AWS::DynamoDB::Table" ||

			d.ResourceType == "AWS::EC2::VPCEndpoint" ||
			d.ResourceType == "AWS::EC2::Instance" ||
			d.ResourceType == "AWS::EC2::VPC" ||
			d.ResourceType == "AWS::EC2::Subnet" ||
			d.ResourceType == "AWS::EC2::NatGateway" ||
			d.ResourceType == "AWS::EC2::InternetGateway" ||
			d.ResourceType == "AWS::EC2::VPCPeeringConnection" ||
			d.ResourceType == "AWS::EC2::SecurityGroup" ||
			d.ResourceType == "AWS::EC2::NetworkInterface" ||
			d.ResourceType == "AWS::EC2::EIP" ||
			d.ResourceType == "AWS::EC2::Volume" ||
			d.ResourceType == "AWS::EC2::NetworkAcl" ||
			d.ResourceType == "AWS::EC2::RouteTable" ||

			d.ResourceType == "AWS::EFS::FileSystem" ||
			d.ResourceType == "AWS::EFS::AccessPoint" ||

			d.ResourceType == "AWS::ElasticLoadBalancingV2::Listener" ||
			d.ResourceType == "AWS::ElasticLoadBalancingV2::LoadBalancer" ||

			d.ResourceType == "AWS::Events::Rule" ||
			d.ResourceType == "AWS::Events::EventBus" ||

			d.ResourceType == "AWS::ECS::Cluster" ||
			d.ResourceType == "AWS::ECS::Service" ||
			d.ResourceType == "AWS::ECS::TaskDefinition" ||

			d.ResourceType == "AWS::IAM::Group" ||
			d.ResourceType == "AWS::IAM::Role" ||
			d.ResourceType == "AWS::IAM::User" ||
			d.ResourceType == "AWS::IAM::Policy" ||

			d.ResourceType == "AWS::KMS::Key" ||
			d.ResourceType == "AWS::Lambda::Function" ||
			d.ResourceType == "AWS::RDS::DBInstance" ||
			d.ResourceType == "AWS::S3::Bucket" ||
			d.ResourceType == "AWS::SNS::Topic" ||
			d.ResourceType == "AWS::ServiceDiscovery::Service" ||
			d.ResourceType == "AWS::ServiceDiscovery::Instance" {
			resuls = append(resuls, d)
		}
	}
	return resuls
}
