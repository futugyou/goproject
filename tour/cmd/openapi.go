package cmd

import (
	"github/go-project/tour/internal/openapi"
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
	},
}

var openapiConfig = openapi.OpenAPIConfig{}

func init() {
	// Priority: flags > .env
	openapiCmd.AddCommand(openapiSubCmd)
	openapiSubCmd.Flags().StringVarP(&openapiConfig.Title, "title", "t", os.Getenv("title"), "openapi name, can also set in .env named 'title'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.Description, "description", "d", os.Getenv("description"), "openapi url, can also set in .env named 'description'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.Version, "version", "r", os.Getenv("version"), "openapi version, can also set in .env named 'version'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.ModelFolder, "folder", "f", os.Getenv("version"), "openapi folder, can also set in .env named 'folder'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.OutputPath, "path", "p", os.Getenv("path"), "openapi output file path, can also set in .env named 'path'")
	openapiSubCmd.Flags().StringVarP(&openapiConfig.OutputType, "type", "", os.Getenv("type"), "openapi output file type, json/yaml, can also set in .env named 'type'")
}
