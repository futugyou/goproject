package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// I want the call to be 'go myplugin subcmd'
// The actual calling method is ./goplugin subcmd
// I don't know how to add the insert plugin to the go command
func main() {
	rootCmd := &cobra.Command{
		Use: "go-myplugin",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "go myplugin",
		},
	}
	subCmd := &cobra.Command{
		Use: "subcmd",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("go myplugin subcmd")
		},
	}
	rootCmd.AddCommand(subCmd)
	rootCmd.Execute()
}
