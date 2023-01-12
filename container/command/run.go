package command

import (
	"fmt"
	"golangproject/container/cgroups/subsystem"
	"golangproject/container/run"

	"github.com/urfave/cli"
)

// 创建namespace隔离的容器进程
// 启动容器
var RunCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "docker volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "docker env",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container args")
		}

		res := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}
		// cmdArray 为容器运行后，执行的第一个命令信息
		// cmdArray[0] 为命令内容, 后面的为命令参数
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		// 要运行的镜像名
		imageName := context.Args().Get(0)
		containerName := context.String("name")
		volume := context.String("v")
		envs := context.StringSlice("e")
		tty := context.Bool("ti")
		detach := context.Bool("d")
		network := context.String("net")
		portMapping := context.StringSlice("p")
		run.Run(cmdArray, tty, detach, res, containerName, imageName, volume, envs, network, portMapping)
		return nil
	},
}
