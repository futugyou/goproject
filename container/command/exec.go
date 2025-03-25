package command

import (
	"golangproject/container/common"
	"golangproject/container/container"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var ExecCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) {
		if os.Getenv(common.EnvExecPid) != "" {
			logrus.Infof("pid callback pid %d", os.Getgid())
			return
		}

		// 我们希望命令格式是docker exec 容器名 命令
		if len(context.Args()) < 2 {
			logrus.Errorf("missing container name or command")
			return
		}

		containerName := context.Args().Get(0)
		var commandArray []string
		commandArray = append(commandArray, context.Args().Tail()...)

		// 执行命令
		if err := container.ExecContainer(containerName, commandArray); err != nil {
			logrus.Errorf("%v", err)
		}
	},
}
