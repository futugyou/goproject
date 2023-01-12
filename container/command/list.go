package command

import (
	"golangproject/container/container"

	"github.com/urfave/cli"
)

var ListCommand = cli.Command{
	Name:  "ps",
	Usage: "list all container",
	Action: func(ctx *cli.Context) error {
		container.ListContainerInfo()
		return nil
	},
}
