package command

import (
	"fmt"
	"golangproject/container/container"

	"github.com/urfave/cli"
)

var RemoveCommand = cli.Command{
	Name:  "rm",
	Usage: "rm a container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing remove container name")
		}
		containerName := ctx.Args().Get(0)
		container.RemoveContainer(containerName)
		return nil
	},
}
