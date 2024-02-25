package cmd

import (
	"github/go-project/tour/internal/mongo2struct"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var mongoLongDesc = strings.Join([]string{
	"generate entity and repository from existed mongodb",
	"command 'g' mean to generate all file",
}, "\n")

var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "mongodb to golang struct and base repository",
	Long:  mongoLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Printf("mongo have no commands named %s, plaese use mongo --help to see detail", strings.Join(args, ","))
		}
	},
}

var mongo2structCmd = &cobra.Command{
	Use:   "g",
	Short: "mongodb to golang struct and base repository",
	Long:  "mongodb to golang struct and base repository",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := mongoDBConfig.ConnectDBDatabase()
		if err != nil {
			log.Println(err)
			return
		}

		m := mongo2struct.NewManager(db, mongoDBConfig.EntityFolder, mongoDBConfig.RepoFolder)
		m.Generator()
	},
}

var mongoDBConfig = mongo2struct.MongoDBConfig{}

func init() {
	mongoCmd.AddCommand(mongo2structCmd)
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.DBName, "dbName", "n", "", "mongodb name, can also set in .env named 'db_name'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.ConnectString, "url", "u", "", "mongodb url, can also set in .env named 'mongodb_url'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.RepoFolder, "repo", "r", "repository", "folder for repository files, can also set in .env named 'repository_folder'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.EntityFolder, "entity", "e", "entity", "folder for entity files, can also set in .env named 'entity_folder'")

	// use godotenv and .env file
	mongodb_url := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	entity_folder := os.Getenv("entity_folder")
	repository_folder := os.Getenv("repository_folder")
	if len(strings.TrimSpace(mongodb_url)) > 0 {
		mongoDBConfig.ConnectString = strings.TrimSpace(mongodb_url)
	}
	if len(strings.TrimSpace(db_name)) > 0 {
		mongoDBConfig.DBName = strings.TrimSpace(db_name)
	}
	if len(strings.TrimSpace(entity_folder)) > 0 {
		mongoDBConfig.EntityFolder = strings.TrimSpace(entity_folder)
	}
	if len(strings.TrimSpace(repository_folder)) > 0 {
		mongoDBConfig.RepoFolder = strings.TrimSpace(repository_folder)
	}
}
