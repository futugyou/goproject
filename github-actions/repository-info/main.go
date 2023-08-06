package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	reponame = flag.String("reponame", "", "repository name, if empty, show all repository.")
	token    = os.Getenv("token")
	ctx      = context.Background()
)

const (
	EOF                = "\r\n"
	multiLineFileDelim = "_GitHubActionsFileCommandDelimeter_"
	multilineFileCmd   = "%s<<" + multiLineFileDelim + EOF + "%s" + EOF + multiLineFileDelim
)

func main() {
	flag.Parse()
	GetGitHubClient()
	GetRepository()

	SetOutput("time", time.Now().Format("2006-01-02 15:04:05"))

	options := make([]ActionResult, len(repos))
	for i := 0; i < len(repos); i++ {
		options[i] = ActionResult{
			Owner:  *repos[i].Owner.Login,
			Name:   *repos[i].Name,
			Branch: *repos[i].DefaultBranch,
		}
	}

	jsonString, _ := json.Marshal(options)
	SetEnv("result", string(jsonString))
}

type ActionResult struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

func SetOutput(key, value string) {
	setGitFile(key, value, "GITHUB_OUTPUT")
}

func SetEnv(key, value string) {
	setGitFile(key, value, "GITHUB_ENV")
}

func setGitFile(key, value, command string) {
	msg := []byte(fmt.Sprintf(multilineFileCmd, key, value) + EOF)

	filepath := os.Getenv(command)
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if _, err := f.Write(msg); err != nil {
		fmt.Println(err)
		return
	}
}
