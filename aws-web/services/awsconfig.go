package services

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"log"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
)

type AwsConfigService struct {
	repository    repository.IAwsConfigRepository
	relRepository repository.IAwsConfigRelationshipRepository
}

func NewAwsConfigService() *AwsConfigService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &AwsConfigService{
		repository:    mongorepo.NewAwsConfigRepository(config),
		relRepository: mongorepo.NewAwsConfigRelationshipRepository(config),
	}
}

func (a *AwsConfigService) SyncFileResources(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	var datas []AwsConfigFileData

	json.Unmarshal(byteValue, &datas)

	if len(datas) == 0 {
		return
	}
	datas = filterResource(datas)
	configs := make([]entity.AwsConfigEntity, 0)

	for _, data := range datas {
		config := createAwsConfigEntity(data)
		configs = append(configs, config)
	}
	log.Println(len(configs))
}

func getId(arn string, resourceID string) string {
	if len(arn) == 0 {
		return resourceID
	}
	return arn
}

func getDataString(con interface{}) string {
	if con == nil {
		return "{}"
	} else {
		d, _ := json.Marshal(con)
		return string(d)
	}
}

// func getVpc(con interface{}) string {

// }
func createAwsConfigEntity(data AwsConfigFileData) entity.AwsConfigEntity {
	config := entity.AwsConfigEntity{
		ID:                           getId(data.ARN, data.ResourceID),
		Label:                        data.ResourceName,
		AccountID:                    data.AwsAccountID,
		Arn:                          data.ARN,
		AvailabilityZone:             data.AvailabilityZone,
		AwsRegion:                    data.AwsRegion,
		Configuration:                getDataString(data.Configuration),
		ConfigurationItemCaptureTime: data.ConfigurationItemCaptureTime,
		ConfigurationItemStatus:      data.ConfigurationItemStatus,
		ConfigurationStateID:         data.ConfigurationStateID,
		ResourceCreationTime:         data.ResourceCreationTime,
		ResourceID:                   data.ResourceID,
		ResourceName:                 data.ResourceName,
		ResourceType:                 data.ResourceType,
		Tags:                         getDataString(data.Tags),
		Version:                      data.ConfigurationItemVersion,
		// VpcID:                        getVpc(data),

		// VpcID                        string      `bson:"vpcId"`
		// SubnetID                     string      `bson:"subnetId"`
		// SubnetIds                    []string    `bson:"subnetIds"`
		// ResourceValue                interface{} `bson:"resourceValue"`
		// State                        interface{} `bson:"state"`
		// Private                      interface{} `bson:"private"`
		// LoggedInURL                  string      `bson:"loggedInURL"`
		// LoginURL                     string      `bson:"loginURL"`
		// Title                        string      `bson:"title"`
		// DBInstanceStatus             string      `bson:"dBInstanceStatus"`
		// Statement                    string      `bson:"statement"`
		// InstanceType                 string      `bson:"instanceType"`
	}
	return config
}

type AwsConfigFileData struct {
	RelatedEvents []string `json:"relatedEvents"`
	Relationships []struct {
		ResourceID   string `json:"resourceId"`
		ResourceType string `json:"resourceType"`
		Name         string `json:"name"`
	} `json:"relationships"`
	Configuration struct {
		CertificateArn          string   `json:"certificateArn"`
		DomainName              string   `json:"domainName"`
		SubjectAlternativeNames []string `json:"subjectAlternativeNames"`
		DomainValidationOptions []struct {
			DomainName       string `json:"domainName"`
			ValidationDomain string `json:"validationDomain"`
			ValidationStatus string `json:"validationStatus"`
			ResourceRecord   struct {
				Name  string `json:"name"`
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"resourceRecord"`
			ValidationMethod string `json:"validationMethod"`
		} `json:"domainValidationOptions"`
		Serial             string    `json:"serial"`
		Subject            string    `json:"subject"`
		Issuer             string    `json:"issuer"`
		CreatedAt          time.Time `json:"createdAt"`
		IssuedAt           time.Time `json:"issuedAt"`
		Status             string    `json:"status"`
		NotBefore          time.Time `json:"notBefore"`
		NotAfter           time.Time `json:"notAfter"`
		KeyAlgorithm       string    `json:"keyAlgorithm"`
		SignatureAlgorithm string    `json:"signatureAlgorithm"`
		InUseBy            []string  `json:"inUseBy"`
		Type               string    `json:"type"`
		RenewalSummary     struct {
			RenewalStatus           string `json:"renewalStatus"`
			DomainValidationOptions []struct {
				DomainName       string `json:"domainName"`
				ValidationDomain string `json:"validationDomain"`
				ValidationStatus string `json:"validationStatus"`
				ResourceRecord   struct {
					Name  string `json:"name"`
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"resourceRecord"`
				ValidationMethod string `json:"validationMethod"`
			} `json:"domainValidationOptions"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"renewalSummary"`
		KeyUsages []struct {
			Name string `json:"name"`
		} `json:"keyUsages"`
		ExtendedKeyUsages []struct {
			Name string `json:"name"`
			OID  string `json:"oID"`
		} `json:"extendedKeyUsages"`
		RenewalEligibility string `json:"renewalEligibility"`
		Options            struct {
			CertificateTransparencyLoggingPreference string `json:"certificateTransparencyLoggingPreference"`
		} `json:"options"`
	} `json:"configuration"`
	SupplementaryConfiguration struct {
		Tags []interface{} `json:"Tags"`
	} `json:"supplementaryConfiguration"`
	Tags struct {
	} `json:"tags"`
	ConfigurationItemVersion     string    `json:"configurationItemVersion"`
	ConfigurationItemCaptureTime time.Time `json:"configurationItemCaptureTime"`
	ConfigurationStateID         int64     `json:"configurationStateId"`
	AwsAccountID                 string    `json:"awsAccountId"`
	ConfigurationItemStatus      string    `json:"configurationItemStatus"`
	ResourceType                 string    `json:"resourceType"`
	ResourceID                   string    `json:"resourceId"`
	ResourceName                 string    `json:"resourceName"`
	ARN                          string    `json:"ARN"`
	AwsRegion                    string    `json:"awsRegion"`
	AvailabilityZone             string    `json:"availabilityZone"`
	ConfigurationStateMd5Hash    string    `json:"configurationStateMd5Hash"`
	ResourceCreationTime         time.Time `json:"resourceCreationTime"`
}

func filterResource(datas []AwsConfigFileData) []AwsConfigFileData {
	resuls := make([]AwsConfigFileData, 0)
	for _, d := range datas {
		if d.ResourceType == "AWS::EC2::VPCEndpoint" ||
			d.ResourceType == "AWS::EC2::VPC" ||
			d.ResourceType == "AWS::ServiceDiscovery::Service" ||
			d.ResourceType == "AWS::Signer::SigningProfile" ||
			d.ResourceType == "AWS::EC2::Subnet" ||
			d.ResourceType == "AWS::AmazonMQ::Broker" ||
			d.ResourceType == "AWS::CloudTrail::Trail" ||
			d.ResourceType == "AWS::EC2::NatGateway" ||
			d.ResourceType == "AWS::EC2::InternetGateway" ||
			d.ResourceType == "AWS::EC2::VPCPeeringConnection" ||
			d.ResourceType == "AWS::EFS::FileSystem" ||
			d.ResourceType == "AWS::IAM::Role" ||
			d.ResourceType == "AWS::RDS::DBInstance" ||
			d.ResourceType == "AWS::SNS::Topic" ||
			d.ResourceType == "AWS::ECS::Cluster" ||
			d.ResourceType == "AWS::IAM::Group" ||
			d.ResourceType == "AWS::ElasticLoadBalancingV2::Listener" ||
			d.ResourceType == "AWS::IAM::User" ||
			d.ResourceType == "AWS::EC2::SecurityGroup" ||
			d.ResourceType == "AWS::EFS::AccessPoint" ||
			d.ResourceType == "AWS::IoT::ProvisioningTemplate" ||
			d.ResourceType == "AWS::EC2::NetworkInterface" ||
			d.ResourceType == "AWS::Route53Resolver::ResolverRuleAssociation" ||
			d.ResourceType == "AWS::RDS::DBSubnetGroup" ||
			d.ResourceType == "AWS::EC2::EIP" ||
			d.ResourceType == "AWS::Redshift::ClusterSubnetGroup" ||
			d.ResourceType == "AWS::ElasticLoadBalancingV2::LoadBalancer" ||
			d.ResourceType == "AWS::ECS::Service" ||
			d.ResourceType == "AWS::EC2::NetworkAcl" ||
			d.ResourceType == "AWS::Lambda::Function" ||
			d.ResourceType == "AWS::S3::Bucket" ||
			d.ResourceType == "AWS::DynamoDB::Table" ||
			d.ResourceType == "AWS::EC2::RouteTable" {
			resuls = append(resuls, d)
		}
	}
	return resuls
}
