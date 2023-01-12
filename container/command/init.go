package command

import (
	"golangproject/container/container"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// 初始化容器内容,挂载proc文件系统，运行用户执行程序
var InitCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		logrus.Infof("init come on")
		return container.RunContainerInitProcess()
	},
}
