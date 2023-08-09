package main

import (
	"fmt"
	"strings"

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
	script.Exec("ls -al").Stdout()
	cloneUrl := fmt.Sprintf("git clone -b %s https://%s@github.com/%s/%s.git .", g.DestBranch, g.DestToken, g.DestOwner, g.DestName)
	result, err := script.Exec(cloneUrl).String()
	if err != nil {
		fmt.Println("CloneDest: ", err)
		return err
	}

	fmt.Println(result)
	return nil
}

func (g *GitCommand) RmoveDest() error {
	scriptstring := "rm -rf *"
	result, err := script.Exec(scriptstring).String()
	if err != nil {
		fmt.Println("RmoveDest: ", err)
		return err
	}

	fmt.Println(result)
	return nil
}

func (g *GitCommand) CloneSource() error {
	script.Exec("mkdir ./sourcerepositoryfolder").Stdout()
	cloneUrl := fmt.Sprintf("git clone -b %s https://%s@github.com/%s/%s.git sourcerepositoryfolder/", g.SourceBranch, g.SourceToken, g.SourceOwner, g.SourceName)
	result, err := script.Exec(cloneUrl).String()
	if err != nil {
		fmt.Println("CloneSource: ", err)
		return err
	}

	fmt.Println(result)
	return nil
}

func (g *GitCommand) CopyToDest() error {
	script.Exec("bash -c 'chmod -R 777 sourcerepositoryfolder'").Stdout()
	script.Exec("bash -c 'rm -rf sourcerepositoryfolder/.git'").Stdout()
	// script.Exec("cp -r sourcerepositoryfolder/.* .").Stdout()
	// result, err := script.Exec("cp -r sourcerepositoryfolder/* .").String()
	script.Exec("bash -c 'mv -f sourcerepositoryfolder/.* .'").Stdout()
	result, err := script.Exec("bash -c 'mv -f sourcerepositoryfolder/* .'").String()
	if err != nil {
		fmt.Println("CopyToDest: ", err)
		return err
	}

	fmt.Println(result)
	return nil
}

func (g *GitCommand) RmoveSourceTemp() error {
	scriptstring := "rm -rf sourcerepositoryfolder/"
	result, err := script.Exec(scriptstring).String()
	if err != nil {
		fmt.Println("RmoveSourceTemp: ", err)
		return err
	}

	fmt.Println(result)
	return nil
}

func (g *GitCommand) GitAdd() {
	scriptstring := "git add ."
	script.Exec(scriptstring).Stdout()
}

func (g *GitCommand) GitStatus() bool {
	status, _ := script.Exec("git status").String()
	fmt.Println(status)
	return !strings.Contains(status, "nothing to commit")
}

func (g *GitCommand) GitCommit() error {
	script.Exec(fmt.Sprintf("git config user.email \"%s\"", botemail)).Stdout()
	script.Exec(fmt.Sprintf("git config user.name  \"%s\"", botuser)).Stdout()

	scriptstring := fmt.Sprintf("git commit -m 'update %s' ", g.DestName)
	result, err := script.Exec(scriptstring).String()
	if err != nil {
		fmt.Println("commit: ", err)
		return err
	}

	fmt.Println(result)
	return nil
}

func (g *GitCommand) GitPush() error {
	scriptstring := fmt.Sprintf("git push https://%s@github.com/%s/%s.git  ", g.DestToken, g.DestOwner, g.DestName)
	result, err := script.Exec(scriptstring).String()
	if err != nil {
		fmt.Println("push: ", err)
		return err
	}

	fmt.Println(result)
	return nil
}

func (g *GitCommand) InitRepository() {
	// script.Exec("echo '' >> ./README.md").Stdout()
	script.Echo("a").WriteFile("./README.md")
	script.Exec("git init").Stdout()
	script.Exec("git add .").Stdout()
	script.Exec(fmt.Sprintf("git config user.email \"%s\"", botemail)).Stdout()
	script.Exec(fmt.Sprintf("git config user.name  \"%s\"", botuser)).Stdout()
	script.Exec("git commit -m \"first commit\"").Stdout()
	script.Exec("git branch -M master").Stdout()
	script.Exec(fmt.Sprintf("git remote add origin https://%s@github.com/%s/%s.git", g.DestToken, g.DestOwner, g.DestName)).Stdout()
	script.Exec("git push -u origin master").Stdout()
	script.Exec("rm -rf ./.git").Stdout()
	script.Exec("rm -rf ./README.md").Stdout()
}
