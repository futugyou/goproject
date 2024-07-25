package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:       "",
	Short:     "",
	Long:      "",
	ValidArgs: []string{"dynamo", "mongo", "openapi", "sql", "time", "word"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute Execute
func Execute() error {
	return rootCmd.Execute()
}
func init() {
	rootCmd.AddCommand(wordCmd)
	rootCmd.AddCommand(timeCmd)
	rootCmd.AddCommand(sqlCmd)
	rootCmd.AddCommand(mongoCmd)
	rootCmd.AddCommand(openapiCmd)
	rootCmd.AddCommand(dynamoCmd)
}
