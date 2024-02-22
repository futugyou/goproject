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
		err := mongoDBConfig.Check()
		if err != nil {
			log.Println(err)
			return
		}

		mongoDBConfig.Generator()
	},
}

var mongoDBConfig = mongo2struct.MongoDBConfig{}

func init() {
	mongoCmd.AddCommand(mongo2structCmd)
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.DBName, "dbName", "n", "", "please input mongodb name !")
	mongo2structCmd.Flags().StringVarP(&mongoDBConfig.ConnectString, "url", "u", "", "please input mongodb url !")
}
