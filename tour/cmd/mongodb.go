package cmd

import (
	"github/go-project/tour/internal/mongo2struct"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "mongodb to golang struct and base repository",
	Long:  "mongodb to golang struct and base repository",
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
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.DBName, "dbName", "n", "", "mongodb name")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.ConnectString, "url", "u", "", "mongodb url")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.RepoFolder, "repo", "r", "repository", "folder for repository files")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.EntityFolder, "entity", "e", "entity", "folder for entity files")
}
