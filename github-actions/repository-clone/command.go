package main

import (
	"fmt"

	"github.com/bitfield/script"
)

const (
	botuser  = "repo-bot-account"
	botemail = "repo-bot-account271538@gmail.com"
)

type GitCommand struct {
	*CloneInfo
}

func NewGitCommand(info *CloneInfo) *GitCommand {
	return &GitCommand{info}
}

func (g *GitCommand) SetConfig() {
	// exec: "cd": executable file not found in $PATH
	// script.Exec("mkdir repositortclonetamp").Stdout()
	// script.Exec("cd repositortclonetamp").Stdout()

	script.Exec("git config --global --add safe.directory '*'").Stdout()
	// fatal: not in a git directory
	// script.Exec(fmt.Sprintf("git config user.email \"%s\"", botemail)).Stdout()
	// fatal: not in a git directory
	// script.Exec(fmt.Sprintf("git config user.name  \"%s\"", botuser)).Stdout()
}

func (g *GitCommand) CloneDest() error {
	cloneUrl := fmt.Sprintf("git clone -b %s https://%s@github.com/%s/%s.git .", g.DestBranch, g.DestToken, g.DestOwner, g.DestName)
	result, err := script.Exec(cloneUrl).String()
	fmt.Println(result, err)
	return err
}

func (g *GitCommand) RmoveDest() error {
	scriptstring := "rm -rf *"
	result, err := script.Exec(scriptstring).String()
	fmt.Println(result, err)
	return err
}

func (g *GitCommand) CloneSource() error {
	script.Exec("mkdir ./sourcerepositoryfolder").Stdout()
	cloneUrl := fmt.Sprintf("git clone -b %s https://%s@github.com/%s/%s.git sourcerepositoryfolder/", g.SourceBranch, g.SourceToken, g.SourceOwner, g.SourceName)
	result, err := script.Exec(cloneUrl).String()
	fmt.Println(result, err)
	return err
}

func (g *GitCommand) CopyToDest() error {
	script.Exec("bash -c 'chmod -R 777 sourcerepositoryfolder'").Stdout()
	script.Exec("bash -c 'rm -rf sourcerepositoryfolder/.git'").Stdout()
	// script.Exec("cp -r sourcerepositoryfolder/.* .").Stdout()
	// result, err := script.Exec("cp -r sourcerepositoryfolder/* .").String()
	script.Exec("bash -c 'mv -f sourcerepositoryfolder/.* .'").Stdout()
	result, err := script.Exec("bash -c 'mv -f sourcerepositoryfolder/* .'").String()
	fmt.Println(result, err)
	return err
}

func (g *GitCommand) RmoveSourceTemp() error {
	scriptstring := "rm -rf sourcerepositoryfolder/"
	result, err := script.Exec(scriptstring).String()
	fmt.Println(result, err)
	return err
}

func (g *GitCommand) GitAdd() error {
	scriptstring := "git add ."
	result, err := script.Exec(scriptstring).String()
	fmt.Println(result, err)
	return err
}

func (g *GitCommand) GitCommit() error {
	script.Exec(fmt.Sprintf("git config user.email \"%s\"", botemail)).Stdout()
	script.Exec(fmt.Sprintf("git config user.name  \"%s\"", botuser)).Stdout()
	scriptstring := fmt.Sprintf("git commit -m 'update %s' ", g.DestName)
	result, err := script.Exec(scriptstring).String()
	fmt.Println(result, err)
	return err
}

func (g *GitCommand) GitPush() error {
	scriptstring := fmt.Sprintf("git push https://%s@github.com/%s/%s.git  ", g.DestToken, g.DestOwner, g.DestName)
	result, err := script.Exec(scriptstring).String()
	fmt.Println("git push: " + result)
	return err
}
