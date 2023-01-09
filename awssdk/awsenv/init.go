package awsenv

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
)

var (
	EmptyContext        context.Context = context.TODO()
	Cfg                 aws.Config
	NamespaceId         string
	CloudMapServiceName string
	UserName            string
	GroupName           string
	// groupPolicyArn      string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	// fmt.Println("AWS_ACCESS_KEY_ID=" + os.Getenv("AWS_ACCESS_KEY_ID"))
	// fmt.Println("AWS_SECRET_ACCESS_KEY=" + os.Getenv("AWS_SECRET_ACCESS_KEY"))
	// fmt.Println("AWS_REGION=" + os.Getenv("AWS_REGION"))

	NamespaceId = os.Getenv("CLOUD_MAP_NAMESPACE")
	CloudMapServiceName = os.Getenv("CLOUD_MAP_SERVICE_NAME")
	UserName = os.Getenv("IAM_USER_NAME")
	GroupName = os.Getenv("IAM_GROUP_NAME")
	// groupPolicyArn = os.Getenv("ATTACHED_GROUP_POLICY_ARN")

	// Load the Shared AWS Configuration (~/.aws/config)
	Cfg, err = config.LoadDefaultConfig(EmptyContext)
	if err != nil {
		log.Fatal(err)
	}

	// internal.CompleteUser(Cfg, GroupName, UserName, groupPolicyArn)
	// Cfg, err = config.LoadDefaultConfig(EmptyContext)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func DeleteAll() {
	// internal.DeleteUser(Cfg, GroupName, UserName, groupPolicyArn)
}
