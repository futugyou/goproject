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
	EmptyContext context.Context = context.TODO()
	Cfg          aws.Config
	NamespaceId  string
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

	// Load the Shared AWS Configuration (~/.aws/config)
	Cfg, err = config.LoadDefaultConfig(EmptyContext)
	if err != nil {
		log.Fatal(err)
	}
}
