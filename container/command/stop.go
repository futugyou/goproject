package command

import (
	"fmt"
	"golangproject/container/container"

	"github.com/urfave/cli"
)

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing stop container name")
		}
		containerName := ctx.Args().Get(0)
		container.StopContainer(containerName)
		return nil
	},
}
