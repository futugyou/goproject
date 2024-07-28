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
	Use:       "openapi",
	Short:     "openapi operate",
	Long:      "openapi operate",
	ValidArgs: []string{"generate", "swag2openapi"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var openapiSubCmdLongDesc = strings.Join([]string{
	"generate an openapi.json or yaml, file spec is:",
	`{
		"spce_version": "3.1.3",
		"title": "this is new title",
		"description": "this is description",
		"apiVersion": "0.0.1",
		"model": "./viewmodel",
		"output": "./openapi.yaml",
		"type": "yaml",
		"apis": [
			{
				"method": "POST",
				"path": "/getall",
				"request": "UserAccountRequest",
				"response": "UserAccountResponse",
				"description": "this is test"
			}
		]
	}`,
}, "\n")

var openapiSubCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate an openapi.json or yaml",
	Long:  openapiSubCmdLongDesc,
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

		astManager := util.NewASTManager(openapiConfig.ModelFolder)

		m, err := openapi.NewManager(*astManager, openapiConfig)
		if err != nil {
			log.Println(err)
			return
		}

		if err = m.GenerateOpenAPI(); err != nil {
			log.Println(err)
		}
	},
}

var swaggerToOpenapiSubCmd = &cobra.Command{
	Use:       "swag2openapi",
	Short:     "convert swagger spec to openapi spec",
	Long:      "convert swagger spec to openapi spec",
	ValidArgs: []string{"yaml", "json"},
	Args:      cobra.MatchAll(cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fileType := ""
		if len(args) == 1 {
			fileType = args[0]
		}
		if err := openapi.ConvertSwaggerToOpenapi(swaggerSpecPath, openapiSpecOutPath, fileType); err != nil {
			log.Println(err)
		}
	},
}

var openapiConfigPath string
var swaggerSpecPath string
var openapiSpecOutPath string

func init() {
	// Priority: flags > .env
	openapiCmd.AddCommand(openapiSubCmd)
	openapiSubCmd.Flags().StringVarP(&openapiConfigPath, "config", "c", "./apiconfig.json", "openapi config file")

	openapiCmd.AddCommand(swaggerToOpenapiSubCmd)
	swaggerToOpenapiSubCmd.Flags().StringVarP(&swaggerSpecPath, "file", "f", "./swagger.yaml", "swagger spec file path")
	swaggerToOpenapiSubCmd.Flags().StringVarP(&openapiSpecOutPath, "output", "o", "./openapi.yaml", "openapi spec file output path")
}
