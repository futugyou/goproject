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
	Short: "generate an openapi.json or yaml",
	Long:  "generate an openapi.json or yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Printf("openapi have no commands named %s, plaese use openapi --help to see detail", strings.Join(args, ","))
		}
	},
}

var openapiSubCmd = &cobra.Command{
	Use:   "generate",
	Short: "openapi to golang struct and base repository",
	Long:  "openapi to golang struct and base repository",
	Run: func(cmd *cobra.Command, args []string) {
		err := openapiConfig.Check()
		if err != nil {
			log.Println(err)
			return
		}
		datas, err := os.ReadFile(openapiConfig.APIRoutePath)
		if err != nil {
			log.Println(err)
			return
		}
		configs := make([]openapi.OperationConfig, 0)
		err = json.Unmarshal(datas, &configs)
		if err != nil {
			log.Println(err)
			return
		}
		l, _ := util.GetStructsFromFolder(openapiConfig.ModelFolder)
		openapiConfig.JsonConfig = configs
		m, err := openapi.NewManager(l, openapiConfig)
		if err != nil {
			log.Println(err)
			return
		}
		err = m.GenerateOpenAPI()
		if err != nil {
			log.Println(err)
		}
	},
}

var openapiConfig = openapi.OpenAPIConfig{}

func init() {
	// Priority: flags > .env
	openapiCmd.AddCommand(openapiSubCmd)
	openapiSubCmd.Flags().StringVarP(&openapiConfig.SpceVersion, "specversion", "v", os.Getenv("spec_version"), "openapi spec version, can also set in .env named 'spec_version'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.Title, "title", "t", os.Getenv("title"), "doc name, can also set in .env named 'title'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.Description, "description", "d", os.Getenv("description"), "doc description, can also set in .env named 'description'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.APIVersion, "api", "i", os.Getenv("version"), "doc version, can also set in .env named 'version'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.ModelFolder, "folder", "f", os.Getenv("version"), "model folder, can also set in .env named 'folder'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.OutputPath, "path", "p", os.Getenv("path"), "output file path, can also set in .env named 'path'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.OutputType, "type", "", os.Getenv("type"), "output file type, json/yaml, can also set in .env named 'type'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.APIRoutePath, "routepath", "", os.Getenv("route_path"), "openapi route file path, can also set in .env named 'route_path'")
}
