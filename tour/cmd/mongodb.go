package cmd

import (
	"github/go-project/tour/internal/mongo2struct"
	"github/go-project/tour/util"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var mongoLongDesc = strings.Join([]string{
	"generate entity and repository from existed mongodb",
	"command 'generate' mean to generate all file",
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
	Use:   "generate",
	Short: "mongodb to golang struct and base repository",
	Long:  "mongodb to golang struct and base repository",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := mongoDBConfig.ConnectDBDatabase()
		if err != nil {
			log.Println(err)
			return
		}

		m := mongo2struct.NewManager(db, mongoDBConfig.EntityFolder, mongoDBConfig.RepoFolder,
			mongoDBConfig.PkgName, mongoDBConfig.CoreFoler, mongoDBConfig.MongoRepoFolder)
		m.Generator()
	},
}

var mongoDBConfig = mongo2struct.MongoDBConfig{}

func init() {
	// Priority: flags > .env
	mongoCmd.AddCommand(mongo2structCmd)
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.DBName, "dbName", "n", os.Getenv("db_name"), "mongodb name, can also set in .env named 'db_name'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.ConnectString, "url", "u", os.Getenv("mongodb_url"), "mongodb url, can also set in .env named 'mongodb_url'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.RepoFolder, "repo", "r", os.Getenv("repository_folder"), "optional. folder for repository files, can also set in .env named 'repository_folder'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.EntityFolder, "entity", "e", os.Getenv("entity_folder"), "optional. folder for entity files, can also set in .env named 'entity_folder'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.PkgName, "package", "p", os.Getenv("package_name"), "optional. package name, can also set in .env named 'package_name'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.CoreFoler, "core", "c", os.Getenv("core_folder"), "optional. folder for core files, can also set in .env named 'core_folder'")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.MongoRepoFolder, "mongo", "m", os.Getenv("mongo_name"), "optional. folder for mongodb repository files, can also set in .env named 'mongo_folder'")

	if len(mongoDBConfig.RepoFolder) == 0 {
		mongoDBConfig.RepoFolder = "repository"
	}
	if len(mongoDBConfig.EntityFolder) == 0 {
		mongoDBConfig.EntityFolder = "entity"
	}
	if len(mongoDBConfig.CoreFoler) == 0 {
		mongoDBConfig.CoreFoler = "core"
	}
	if len(mongoDBConfig.MongoRepoFolder) == 0 {
		mongoDBConfig.MongoRepoFolder = "mongorepo"
	}
	if len(mongoDBConfig.PkgName) == 0 {
		mongoDBConfig.PkgName = util.GetModuleName()
	}
}
