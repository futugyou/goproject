package cmd

import (
	"encoding/json"
	"github/go-project/tour/internal/openapi"
	"github/go-project/tour/util"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "openapi generate",
	Long:  "openapi generate",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Printf("openapi have no commands named %s, plaese use openapi --help to see detail", strings.Join(args, ","))
		}
	},
}

var openapiSubCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate an openapi.json or yaml",
	Long:  "generate an openapi.json or yaml",
	Run: func(cmd *cobra.Command, args []string) {
		datas, err := os.ReadFile(openapiConfigPath)
		if err != nil {
			log.Println(err)
			return
		}

		var openapiConfig = openapi.OpenAPIConfig{}
		if err = json.Unmarshal(datas, &openapiConfig); err != nil {
			log.Println(err)
			return
		}

		if err = openapiConfig.Check(); err != nil {
			log.Println(err)
			return
		}

		structInfos, err := util.GetStructsFromFolder(openapiConfig.ModelFolder)
		if err != nil {
			log.Println(err)
			return
		}

		m, err := openapi.NewManager(structInfos, openapiConfig)
		if err != nil {
			log.Println(err)
			return
		}

		if err = m.GenerateOpenAPI(); err != nil {
			log.Println(err)
		}
	},
}

var openapiConfigPath string

func init() {
	// Priority: flags > .env
	openapiCmd.AddCommand(openapiSubCmd)
	openapiSubCmd.Flags().StringVarP(&openapiConfigPath, "config", "c", "./apiconfig.json", "openapi config file")
}
