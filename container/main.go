package main

import (
	"golangproject/container/command"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `go-docker`

func main() {
	app := cli.NewApp()
	app.Name = "go-docker"
	app.Usage = usage

	app.Commands = []cli.Command{
		command.InitCommand,
		command.RunCommand,
		command.LogCommand,
		command.ListCommand,
		command.StopCommand,
		command.RemoveCommand,
		command.CommitCommand,
		command.NetworkCommand,
	}
	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
