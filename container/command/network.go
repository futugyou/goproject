package command

import (
	"fmt"
	"golangproject/container/network"

	"github.com/urfave/cli"
)

var NetworkCommand = cli.Command{
	Name:  "network",
	Usage: "container network commands",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}

				err := network.Init()
				if err != nil {
					return err
				}

				driver := context.String("driver")
				subnet := context.String("subnet")
				name := context.Args()[0]
				if err := network.CreateNetwork(driver, subnet, name); err != nil {
					return err
				}
				return nil
			},
		},

		{
			Name:  "list",
			Usage: "list container network",
			Action: func(context *cli.Context) error {
				if err := network.Init(); err != nil {
					return err
				}
				return network.ListNetwork()
			},
		},

		{
			Name:  "remove",
			Usage: "remove container network",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}
				if err := network.Init(); err != nil {
					return err
				}
				if err := network.DeleteNetwork(context.Args()[0]); err != nil {
					return err
				}
				return nil
			},
		},
	},
}
