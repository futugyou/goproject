package cmd

import (
	"github/go-project/tour/internal/dynamo2struct"
	"github/go-project/tour/util"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var dynamoLongDesc = strings.Join([]string{
	"generate entity and repository from existed dynamodb",
	"command 'generate' mean to generate all file",
}, "\n")

var dynamoCmd = &cobra.Command{
	Use:   "dynamo",
	Short: "dynamodb to golang struct and base repository",
	Long:  dynamoLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Printf("dynamo have no commands named %s, plaese use dynamo --help to see detail", strings.Join(args, ","))
		}
	},
}

var dynamo2structCmd = &cobra.Command{
	Use:   "generate",
	Short: "dynamodb to golang struct and base repository",
	Long:  "dynamodb to golang struct and base repository",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var dynamoDBConfig = dynamo2struct.DynamoDBConfig{}

func init() {
	// Priority: flags > .env
	dynamoCmd.AddCommand(dynamo2structCmd)
	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.AccessKey, "key", "k", os.Getenv("AWS_ACCESS_KEY_ID"), "dynamodb AccessKey, can also set in .env named 'AWS_ACCESS_KEY_ID'")
	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.AccessSecret, "secret", "s", os.Getenv("AWS_SECRET_ACCESS_KEY"), "dynamodb AccessSecret, can also set in .env named 'AWS_SECRET_ACCESS_KEY'")
	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.Region, "region", "g", os.Getenv("AWS_REGION"), "dynamodb Region, can also set in .env named 'AWS_REGION'")

	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.RepoFolder, "repo", "r", os.Getenv("repository_folder"), "optional. folder for repository files, can also set in .env named 'repository_folder'")
	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.EntityFolder, "entity", "e", os.Getenv("entity_folder"), "optional. folder for entity files, can also set in .env named 'entity_folder'")
	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.PkgName, "package", "p", os.Getenv("package_name"), "optional. package name, can also set in .env named 'package_name'")
	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.CoreFoler, "core", "c", os.Getenv("core_folder"), "optional. folder for core files, can also set in .env named 'core_folder'")
	dynamo2structCmd.Flags().StringVarP(&dynamoDBConfig.DynamoRepoFolder, "dynamo", "m", os.Getenv("dynamo_name"), "optional. folder for dynamodb repository files, can also set in .env named 'dynamo_folder'")

	if len(dynamoDBConfig.RepoFolder) == 0 {
		dynamoDBConfig.RepoFolder = "repository"
	}
	if len(dynamoDBConfig.EntityFolder) == 0 {
		dynamoDBConfig.EntityFolder = "entity"
	}
	if len(dynamoDBConfig.CoreFoler) == 0 {
		dynamoDBConfig.CoreFoler = "core"
	}
	if len(dynamoDBConfig.DynamoRepoFolder) == 0 {
		dynamoDBConfig.DynamoRepoFolder = "dynamorepo"
	}
	if len(dynamoDBConfig.PkgName) == 0 {
		dynamoDBConfig.PkgName = util.GetModuleName()
	}
}
